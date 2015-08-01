package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
)

const predictedUsersCount = 1500
const predictedClastersCount = 10
const predictedClasterSize = 6

var userEmails []string
var connections map[int]map[int]struct{}
var clusters []map[int]struct{}

func main() {
	readInput()
	cleanConnections()
	findClusters()
	printClusters()

	/*fmt.Printf("users %d, active users %d, clusters %d\n", len(userEmails), len(connections), len(clusters))

	file, _ := os.Open(os.Args[1])
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}*/
}

func findClusters() {
	clusters = make([]map[int]struct{}, 0, predictedClastersCount)

	// Iterate over all connections
	for userId1, user1Connections := range connections {
		for userId2, _ := range user1Connections {
			if userId2 < userId1 {
				continue // Skip mirrored connections
			} else {
				findSuperClasters(map[int]struct{}{userId1: struct{}{}, userId2: struct{}{}})
			}
		}
	}
}

func findSuperClasters(baseCluster map[int]struct{}) {
	hasSuperClusters := false

	// Iterate over all users and try to add them to base cluster
	for potentialUserId, potentialUserConnections := range connections {
		// Skip users already in base cluster
		_, present := baseCluster[potentialUserId]
		if present {
			continue
		}

		// Check if potential user connected to cluster
		var connected bool
		for clusterUserId, _ := range baseCluster {
			_, connected = potentialUserConnections[clusterUserId]
			if !connected { // Not connected to one of cluster users
				break
			}
		}

		if connected {
			hasSuperClusters = true

			// Recursively find super clusters of the cluster with added new user
			biggerBaseCluster := createClusterCopy(baseCluster)
			biggerBaseCluster[potentialUserId] = struct{}{}
			findSuperClasters(biggerBaseCluster)
		}
	}

	if (!hasSuperClusters) && (len(baseCluster) >= 3) && (!isClusterExists(baseCluster)) {
		// Add final cluster to list of all clusters
		clusters = append(clusters, baseCluster)
	}
}

func printClusters() {
	clusterStrings := make([]string, len(clusters))

	for clusterIndex, cluster := range clusters {
		clusterEmails := make([]string, 0, predictedClasterSize)
		for userId, _ := range cluster {
			clusterEmails = append(clusterEmails, userEmails[userId])
		}

		sort.Strings(clusterEmails)

		var b bytes.Buffer
		for index, email := range clusterEmails {
			if index != 0 {
				b.WriteString(", ")
			}
			b.WriteString(email)
		}

		clusterStrings[clusterIndex] = b.String()
	}

	sort.Strings(clusterStrings)

	for _, clusterString := range clusterStrings {
		fmt.Println(clusterString)
	}
}

func createClusterCopy(cluster map[int]struct{}) map[int]struct{} {
	clusterCopy := make(map[int]struct{})
	for userId, _ := range cluster {
		clusterCopy[userId] = struct{}{}
	}

	return clusterCopy
}

func isClusterExists(clusterToCheck map[int]struct{}) bool {
	for _, cluster := range clusters {
		if len(clusterToCheck) != len(cluster) {
			continue
		}

		var present bool
		for userId, _ := range clusterToCheck {
			_, present = cluster[userId]
			if !present {
				break
			}
		}

		if present {
			return true
		}
	}

	return false
}

func readInput() {
	userEmails = make([]string, 0, predictedUsersCount)
	connections = make(map[int]map[int]struct{})

	f, _ := os.Open(os.Args[1])
	defer f.Close()
	scanner := bufio.NewScanner(bufio.NewReader(f))
	scanner.Split(bufio.ScanWords)

	userIds := make(map[string]int)

	for scanner.Scan() {
		// Skip 6 words
		for i := 0; i < 5; i++ {
			scanner.Scan()
		}

		// Read 2 emails, define user id by email, create new users if needed

		scanner.Scan()
		email1 := scanner.Text()

		userId1, present := userIds[email1]
		if !present {
			userId1 = len(userEmails)
			userIds[email1] = userId1
			userEmails = append(userEmails, email1)
			connections[userId1] = make(map[int]struct{})
		}

		scanner.Scan()
		email2 := scanner.Text()

		userId2, present := userIds[email2]
		if !present {
			userId2 = len(userEmails)
			userIds[email2] = userId2
			userEmails = append(userEmails, email2)
			connections[userId2] = make(map[int]struct{})
		}

		// Store connection
		connections[userId1][userId2] = struct{}{}
	}
}

func cleanConnections() {
	for userId1, user1Connections := range connections {
		// Remove single side connections
		for userId2, _ := range user1Connections {
			_, present := connections[userId2][userId1]
			if !present {
				delete(user1Connections, userId2)
			}
		}

		// Remove users without connections
		if len(user1Connections) == 0 {
			delete(connections, userId1)
		}
	}
}
