package manifest

type Action string

const CreateAction = "Create"
const UpdateAction = "Update"
const DeleteAction = "Delete"

// Message wraps a manifest with information related to the Loader that loaded it
type Message struct {
	// Manifest loaded
	Manifest *Manifest

	// Action to perform (e.g. create, update, delete...)
	Action Action
}
