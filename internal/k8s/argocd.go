package k8s

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
)

// VerifyArgoCDReadiness waits for critical resources within ArgoCD to be ready
// and only returns once they're all healthy
//
// This helps prevent race conditions and timeouts
func VerifyArgoCDReadiness(clientset *kubernetes.Clientset, highAvailabilityEnabled bool) (bool, error) {
	// Wait for ArgoCD StatefulSet Pods to transition to Running
	argoCDStatefulSet, err := ReturnStatefulSetObject(
		clientset,
		"app.kubernetes.io/part-of",
		"argocd",
		"argocd",
		240,
	)
	if err != nil {
		return false, fmt.Errorf("error finding ArgoCD Application Controller StatefulSet: %s", err)
	}
	_, err = WaitForStatefulSetReady(clientset, argoCDStatefulSet, 240, false)
	if err != nil {
		return false, fmt.Errorf("error waiting for ArgoCD Application Controller StatefulSet ready state: %s", err)
	}

	// argocd-server Deployment
	argoCDServerDeployment, err := ReturnDeploymentObject(
		clientset,
		"app.kubernetes.io/name",
		"argocd-server",
		"argocd",
		240,
	)
	if err != nil {
		log.Info().Msgf("Error finding ArgoCD server deployment: %s", err)
	}
	_, err = WaitForDeploymentReady(clientset, argoCDServerDeployment, 240)
	if err != nil {
		log.Info().Msgf("Error waiting for ArgoCD server deployment ready state: %s", err)
	}

	// Wait for additional ArgoCD Pods to transition to Running
	// This is related to a condition where apps attempt to deploy before
	// repo, redis, or other health checks are passing
	//
	// This can cause future steps to break since the registry app
	// may never apply

	// argocd-repo-server
	argoCDRepoDeployment, err := ReturnDeploymentObject(
		clientset,
		"app.kubernetes.io/name",
		"argocd-repo-server",
		"argocd",
		240,
	)
	if err != nil {
		return false, fmt.Errorf("error finding ArgoCD repo deployment: %s", err.Error())
	}
	_, err = WaitForDeploymentReady(clientset, argoCDRepoDeployment, 240)
	if err != nil {
		return false, fmt.Errorf("error waiting for ArgoCD repo deployment ready state: %s", err.Error())
	}

	// high availability components
	if highAvailabilityEnabled {
		// argocd-redis-ha-haproxy Deployment
		argoCDRedisHAhaproxyDeployment, err := ReturnDeploymentObject(
			clientset,
			"app.kubernetes.io/name",
			"argocd-redis-ha-haproxy",
			"argocd",
			240,
		)
		if err != nil {
			return false, fmt.Errorf("error finding ArgoCD argocd-redis-ha-haproxy Deployment: %s", err)
		}
		_, err = WaitForDeploymentReady(clientset, argoCDRedisHAhaproxyDeployment, 240)
		if err != nil {
			return false, fmt.Errorf("error waiting for ArgoCD argocd-redis-ha-haproxy deployment ready state: %s", err.Error())
		}

		// argocd-redis-ha StatefulSet
		argoCDRedisHAServerStatefulSet, err := ReturnStatefulSetObject(
			clientset,
			"app.kubernetes.io/name",
			"argocd-redis-ha",
			"argocd",
			240,
		)
		if err != nil {
			return false, fmt.Errorf("error finding ArgoCD argocd-redis-ha StatefulSet: %s", err.Error())
		}
		_, err = WaitForStatefulSetReady(clientset, argoCDRedisHAServerStatefulSet, 240, false)
		if err != nil {
			return false, fmt.Errorf("error waiting for ArgoCD argocd-redis-ha StatefulSet ready state: %s", err.Error())
		}
	} else {
		// non-high availability components
		// argocd-redis Deployment
		argoCDRedisDeployment, err := ReturnDeploymentObject(
			clientset,
			"app.kubernetes.io/name",
			"argocd-redis",
			"argocd",
			240,
		)
		if err != nil {
			return false, fmt.Errorf("error finding ArgoCD argocd-redis Deployment: %s", err)
		}
		_, err = WaitForDeploymentReady(clientset, argoCDRedisDeployment, 240)
		if err != nil {
			return false, fmt.Errorf("error waiting for ArgoCD argocd-redis Deployment ready state: %s", err)
		}
	}

	return true, nil
}
