BASEUPLOADURL="https://nostr.build/api/v2/upload/files"
NSEC=$(cat config.json|jq -r .nsec)

# Auth to nostr build
AUTHJSON=$(nak event --kind 27235 --tag u=uploadUrl --tag method=uploadMethod --sec $NSEC)
# todo: stringify, then convert to base64 and prepare auth header as follows
#         let jsonAuthEvent = JSON.stringify(signedAuthEvent);
#         let base64AuthEvent = btoa(jsonAuthEvent);
#         authHeader = `Nostr: ${base64AuthEvent}`;