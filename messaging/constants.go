package messaging

const CONSUMING_QUEUE string = "reverseproxy-to-admin"
const PUBLISHING_QUEUE string = "admin-to-reverseproxy"

// for publishing
const ADD_REPLICA string = "add-replica"
const REMOVE_REPLICA string = "remove-replica"
const NEW_PARAMETERS string = "new-parameters"

// for consuming
const ADDED_REPLICA string = "replica-added"
const REMOVED_REPLICA string = "replica-removed"
const STATISTICS string = "statistics"
const PARAMETERS_UPDATED = "parameters-updated"
const PARAMETERS_UPDATE_FAILED = "parameters-update-failed"
