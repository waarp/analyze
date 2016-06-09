package waarp

import (
	"fmt"
	"strings"
)

type ErrorCode string

const (
	EC_WM_UNKNOWN             ErrorCode = ""
	EC_INIT_OK                ErrorCode = "i"
	EC_PRE_PROCESSING_OK      ErrorCode = "B"
	EC_TRANSFER_OK            ErrorCode = "X"
	EC_POST_PROCESSING_OK     ErrorCode = "P"
	EC_COMPLETE_OK            ErrorCode = "O"
	EC_CONNECTION_IMPOSSIBLE  ErrorCode = "C"
	EC_SERVER_OVERLOADED      ErrorCode = "l"
	EC_BAD_AUTHENT            ErrorCode = "A"
	EC_EXTERNAL_OP            ErrorCode = "E"
	EC_TRANSFER_ERROR         ErrorCode = "T"
	EC_MD5_ERROR              ErrorCode = "M"
	EC_DISCONNECTION          ErrorCode = "D"
	EC_REMOTE_SHUTDOWN        ErrorCode = "r"
	EC_FINAL_OP               ErrorCode = "F"
	EC_UNIMPLEMENTED          ErrorCode = "U"
	EC_SHUTDOWN               ErrorCode = "S"
	EC_REMOTE_ERROR           ErrorCode = "R"
	EC_INTERNAL               ErrorCode = "I"
	EC_STOPPED_TRANSFER       ErrorCode = "H"
	EC_CANCELED_TRANSFER      ErrorCode = "K"
	EC_WARNING                ErrorCode = "W"
	EC_UNKNOWN                ErrorCode = "-"
	EC_QUERY_ALREADY_FINISHED ErrorCode = "Q"
	EC_QUERY_STILL_RUNNING    ErrorCode = "s"
	EC_NOT_KNOWN_HOST         ErrorCode = "N"
	EC_QUERY_REMOTELY_UNKNOWN ErrorCode = "u"
	EC_FILE_NOT_FOUND         ErrorCode = "f"
	EC_COMMAND_NOT_FOUND      ErrorCode = "c"
	EC_PASS_THROUGH_MODE      ErrorCode = "p"
	EC_RUNNING                ErrorCode = "z"
	EC_INCORRECT_COMMAND      ErrorCode = "n"
	EC_FILE_NOT_ALLOWED       ErrorCode = "a"
	EC_SIZE_NOT_ALLOWED       ErrorCode = "d"
	// EC_LOOP_SELF_REQUESTED_HOST ErrorCode = "N"
)

var errorCodeMessages = map[ErrorCode]string{
	EC_WM_UNKNOWN:             "Erreur inconnue de Waarp Manager",
	EC_INIT_OK:                "Transfert initialisé avec succès",
	EC_PRE_PROCESSING_OK:      "Pré-traitements efféctués avec succès",
	EC_TRANSFER_OK:            "Données transferées avec succès",
	EC_POST_PROCESSING_OK:     "Post-traitements efféctués avec succès",
	EC_COMPLETE_OK:            "Transfert terminé avec succès",
	EC_CONNECTION_IMPOSSIBLE:  "La connexion avec le partenaire a échouée",
	EC_SERVER_OVERLOADED:      "L'action demandée est ne peut être effectuée par manque de ressources",
	EC_BAD_AUTHENT:            "L'authentification entre les partenaires a échouée",
	EC_EXTERNAL_OP:            "Une opération externe a échouée",
	EC_TRANSFER_ERROR:         "Une erreur s'est produite durant le transfert des données",
	EC_MD5_ERROR:              "La vérification MD5 du transfert a échouée",
	EC_DISCONNECTION:          "Le réseau n'est plus accessible",
	EC_REMOTE_SHUTDOWN:        "L'arrêt du serveur a été demandé",
	EC_FINAL_OP:               "Une erreur s'est produite durant la finalisation du transfert",
	EC_UNIMPLEMENTED:          "Une fonctionnalité demandée n'est pas implémentée",
	EC_SHUTDOWN:               "L'arrêt du serveur est en cours",
	EC_REMOTE_ERROR:           "Une erreur s'est produite sur le partenaire distant",
	EC_INTERNAL:               "Une erreur interne s'est produite",
	EC_STOPPED_TRANSFER:       "Le transfert a été stoppé par un administrateur",
	EC_CANCELED_TRANSFER:      "Le transfert a été annulé par un administrateur",
	EC_WARNING:                "Une opération externe a renvoyé un warning",
	EC_UNKNOWN:                "Une erreur inconnue s'est produite",
	EC_QUERY_ALREADY_FINISHED: "La requête est déjà terminée chez le partenaire",
	EC_QUERY_STILL_RUNNING:    "La requête est toujours en cours",
	EC_NOT_KNOWN_HOST:         "Le partenaire demandé est inconnu",
	// EC_LOOP_SELF_REQUESTED_HOST: "Auto-transfert impossible",
	EC_QUERY_REMOTELY_UNKNOWN: "La reqête n'a pas été trouvée sur le partenaire distant",
	EC_FILE_NOT_FOUND:         "Le fichier a transférer n'a pas été trouvé",
	EC_COMMAND_NOT_FOUND:      "La commande demandée n'a pas été trouvée",
	EC_PASS_THROUGH_MODE:      "Le transfert est en mode passthrough, et l'action demandée est incompatible avec ce mode",
	EC_RUNNING:                "L'étape est en cours",
	EC_INCORRECT_COMMAND:      "La commande demandée est incorrecte",
	EC_FILE_NOT_ALLOWED:       "Ce fichier n'est pas autorisé",
	EC_SIZE_NOT_ALLOWED:       "La taille n'est pas autorisée",
}

var errorCodeNames = map[ErrorCode]string{
	EC_WM_UNKNOWN:             "EC_WM_UNKNOWN",
	EC_INIT_OK:                "EC_INIT_OK",
	EC_PRE_PROCESSING_OK:      "EC_PRE_PROCESSING_OK",
	EC_TRANSFER_OK:            "EC_TRANSFER_OK",
	EC_POST_PROCESSING_OK:     "EC_POST_PROCESSING_OK",
	EC_COMPLETE_OK:            "EC_COMPLETE_OK",
	EC_CONNECTION_IMPOSSIBLE:  "EC_CONNECTION_IMPOSSIBLE",
	EC_SERVER_OVERLOADED:      "EC_SERVER_OVERLOADED",
	EC_BAD_AUTHENT:            "EC_BAD_AUTHENT",
	EC_EXTERNAL_OP:            "EC_EXTERNAL_OP",
	EC_TRANSFER_ERROR:         "EC_TRANSFER_ERROR",
	EC_MD5_ERROR:              "EC_MD5_ERROR",
	EC_DISCONNECTION:          "EC_DISCONNECTION",
	EC_REMOTE_SHUTDOWN:        "EC_REMOTE_SHUTDOWN",
	EC_FINAL_OP:               "EC_FINAL_OP",
	EC_UNIMPLEMENTED:          "EC_UNIMPLEMENTED",
	EC_SHUTDOWN:               "EC_SHUTDOWN",
	EC_REMOTE_ERROR:           "EC_REMOTE_ERROR",
	EC_INTERNAL:               "EC_INTERNAL",
	EC_STOPPED_TRANSFER:       "EC_STOPPED_TRANSFER",
	EC_CANCELED_TRANSFER:      "EC_CANCELED_TRANSFER",
	EC_WARNING:                "EC_WARNING",
	EC_UNKNOWN:                "EC_UNKNOWN",
	EC_QUERY_ALREADY_FINISHED: "EC_QUERY_ALREADY_FINISHED",
	EC_QUERY_STILL_RUNNING:    "EC_QUERY_STILL_RUNNING",
	EC_NOT_KNOWN_HOST:         "EC_NOT_KNOWN_HOST",
	EC_QUERY_REMOTELY_UNKNOWN: "EC_QUERY_REMOTELY_UNKNOWN",
	EC_FILE_NOT_FOUND:         "EC_FILE_NOT_FOUND",
	EC_COMMAND_NOT_FOUND:      "EC_COMMAND_NOT_FOUND",
	EC_PASS_THROUGH_MODE:      "EC_PASS_THROUGH_MODE",
	EC_RUNNING:                "EC_RUNNING",
	EC_INCORRECT_COMMAND:      "EC_INCORRECT_COMMAND",
	EC_FILE_NOT_ALLOWED:       "EC_FILE_NOT_ALLOWED",
	EC_SIZE_NOT_ALLOWED:       "EC_SIZE_NOT_ALLOWED",
}

func (ec ErrorCode) String() string {
	return errorCodeMessages[ec]
}

func (ec ErrorCode) FullString() string {
	return fmt.Sprintf("%s: %s", errorCodeNames[ec], errorCodeMessages[ec])
}

// JSON Serialization
func (ec *ErrorCode) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ec.FullString() + `"`), nil
}

// JSON Deserialization
func (ec *ErrorCode) UnmarshalJSON(val []byte) error {
	value := ErrorCode(strings.TrimSpace(string(val[1 : len(val)-1])))
	if _, ok := errorCodeNames[value]; ok {
		*ec = value
	} else {
		*ec = EC_WM_UNKNOWN
	}
	return nil
}
