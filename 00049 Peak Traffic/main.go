package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const predictedUsersCount = 1500
const predictedClastersCount = 10

var userEmails []string
var connections map[int]map[int]struct{}
var clusters []map[int]struct{}

func main() {
	readInput()
	cleanConnections()
	findClusters()
	printClusters()
	//fmt.Printf("users %d, active users %d, clusters %d\n", len(userEmails), len(connections), len(clusters))
}

func findClusters() {
	clusters = make([]map[int]struct{}, 0, predictedClastersCount)

	// Iterate over all connections
	for userId1, user1Connections := range connections {
		for userId2, _ := range user1Connections {
			if userId2 < userId1 {
				continue // Skip mirrored connections
			} else if !isClusterExists(userId1, userId2) {
				doCluster(userId1, userId2)
			}
		}
	}
}

func doCluster(userId1 int, userId2 int) {
	cluster := map[int]struct{}{userId1: struct{}{}, userId2: struct{}{}}

	// Iterate over all users and try to add them to cluster
	for potentialUserId, potentialUserConnections := range connections {
		if (potentialUserId == userId1) || (potentialUserId == userId2) { // Skip users already in cluster
			continue
		}

		// Check if potential user connected to cluster
		var connected bool
		for clusterUserId, _ := range cluster {
			if _, connected = potentialUserConnections[clusterUserId]; !connected { // Not connected to one of cluster users
				break
			}
		}

		if connected {
			cluster[potentialUserId] = struct{}{} // Add user to cluster
		}
	}

	if len(cluster) >= 3 {
		clusters = append(clusters, cluster)
	}
}

func printClusters() {
	clusterStrings := make([]string, len(clusters))

	for clusterIndex, cluster := range clusters {
		clusterEmails := make([]string, 0, len(cluster))
		for userId, _ := range cluster {
			clusterEmails = append(clusterEmails, userEmails[userId])
		}

		sort.Strings(clusterEmails)
		clusterStrings[clusterIndex] = strings.Join(clusterEmails, ", ")
	}

	sort.Strings(clusterStrings)

	for _, clusterString := range clusterStrings {
		fmt.Println(clusterString)
	}
}

func isClusterExists(userId1 int, userId2 int) bool {
	for _, cluster := range clusters {
		if _, present := cluster[userId1]; present {
			if _, present = cluster[userId2]; present {
				return true
			}
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

		// Read 2 emails, define user ID by email, create new users if needed
		var connectionUserIds [2]int
		for i := 0; i < len(connectionUserIds); i++ {
			scanner.Scan()
			email := scanner.Text()

			userId, present := userIds[email]
			if !present {
				userId = len(userEmails)
				userIds[email] = userId
				userEmails = append(userEmails, email)
				connections[userId] = make(map[int]struct{})
			}

			connectionUserIds[i] = userId
		}

		// Store connection
		connections[connectionUserIds[0]][connectionUserIds[1]] = struct{}{}
	}
}

func cleanConnections() {
	for userId1, user1Connections := range connections {
		// Remove single side connections
		for userId2, _ := range user1Connections {
			if _, present := connections[userId2][userId1]; !present {
				delete(user1Connections, userId2)
			}
		}

		// Remove users without connections
		if len(user1Connections) == 0 {
			delete(connections, userId1)
		}
	}
}
