package server

import (
	"assignment/misc"
	_ "assignment/misc"
	"assignment/repository"
	_ "assignment/repository"
	"encoding/json"
	"net/http"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(writer http.ResponseWriter, req *http.Request) {
	// Invalid req method
	if req.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Invalid req body
	var userCred UserCredentials
	err := json.NewDecoder(req.Body).Decode(&userCred)
	if err != nil {
		http.Error(writer, "Error decoding body", http.StatusBadRequest)
		return
	}
	// Validate username and pass
	userFromDB, err := repository.GetUserFromDB(userCred.Username)
	if err == nil && userCred.Password == userFromDB.Password {

		// Create token (json string) with - username, pass, accessLevel
		tokenMap := map[string]interface{}{
			"username":     userFromDB.Username,
			"password":     userFromDB.Password,
			"access_level": string(userFromDB.AccessLevel),
			"entity_id":    userFromDB.EntityId,
		}
		tokenJsonString, err := json.Marshal(tokenMap)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		// Create json response
		response := map[string]string{
			"token": string(tokenJsonString),
		}
		jsonEncoder := json.NewEncoder(writer)
		jsonEncoder.Encode(response)
	} else {
		http.Error(writer, "Invalid username or password", http.StatusBadRequest)
		return
	}
}

func PersonHandler(writer http.ResponseWriter, req *http.Request) {

	// Validate token
	headers := req.Header
	tokenJsonString := headers["Token"]
	if tokenJsonString == nil {
		http.Error(writer, "Invalid token", http.StatusUnauthorized)
		return
	}
	userDetails := make(map[string]interface{})
	token := tokenJsonString[0]
	// Get data from token
	tokenWithoutSlashes := misc.RemoveSlashes(token)
	json.Unmarshal([]byte(tokenWithoutSlashes), &userDetails)

	username := userDetails["username"]
	accessLevel := userDetails["access_level"]
	password := userDetails["password"]
	entityId := userDetails["entity_id"]
	if accessLevel == nil || username == nil || password == nil || entityId == nil {
		http.Error(writer, "Invalid token", http.StatusUnauthorized)
		return
	}
	if accessLevel != "ADMIN" {
		http.Error(writer, "User is not an admin", http.StatusUnauthorized)
		return
	}

	// After validation, If POST API, call addPersonDetails()
	if req.Method == http.MethodPost {
		addPersonDetails(writer, req)
		return
	}
	// GET API
	// Fetch details from details table
	if entityIdInt, ok := entityId.(float64); ok {
		userDetailsFromDB, err := repository.GetUserDetailsFromDB(entityIdInt)

		if err != nil {
			http.Error(writer, "Failed to fetch row with given userId - "+err.Error(), http.StatusInternalServerError)
			return
		} else {
			response := map[string]interface{}{
				"first_name":    userDetailsFromDB.FirstName,
				"last_name":     userDetailsFromDB.LastName,
				"age":           userDetailsFromDB.Age,
				"address":       userDetailsFromDB.Address,
				"gender":        userDetailsFromDB.Gender,
				"email":         userDetailsFromDB.Email,
				"mobile_number": userDetailsFromDB.MobileNumber,
			}
			writer.WriteHeader(http.StatusOK)
			writer.Header().Set("Content-Type", "application/json")
			jsonEncoder := json.NewEncoder(writer)
			err := jsonEncoder.Encode(response)
			if err != nil {
				http.Error(writer, "Failed to encode"+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		http.Error(writer, "Cannot fetch the user atm", http.StatusInternalServerError)
		return
	}
}

func addPersonDetails(writer http.ResponseWriter, req *http.Request) {

	var userDetails misc.UserDetailsRequest
	err := json.NewDecoder(req.Body).Decode(&userDetails)
	if err != nil {
		http.Error(writer, "Error decoding body", http.StatusBadRequest)
		return
	}
	// Add data to table
	err = repository.AddPerson(userDetails)
	writer.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Successfully added new data",
	}
	jsonEncoder := json.NewEncoder(writer)
	if err != nil {
		response["message"] = err.Error()
	}

	jsonEncoder.Encode(response)
}
