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
const predictedClasterSize = 6

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
			_, connected = potentialUserConnections[clusterUserId]
			if !connected { // Not connected to one of cluster users
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
		clusterEmails := make([]string, 0, predictedClasterSize)
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
		_, present := cluster[userId1]

		if present {
			_, present = cluster[userId2]

			if present {
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
