# Uploader

Script to bulk upload images to a media storage location (using nostr.build), and with resulting URL, create nostr events

Assume nak (https://github.com/fiatjaf/nak) is installed and on path

# Config

In config.json

```json
{
    "nsec": "nsec value to use for signing nostr events",
    "imagesource": "/path/to/images",
    "category": "#pepe",
}

```

# Nostr.Build upload... based on the javascript equivalent

uploadUrl = 'https://nostr.build/api/v2/upload/files';


 const uploadMethod = 'POST';
    let doAuth = ((localStorage.getItem(`fileUpload.auth`) ?? 'false') == 'true');
    let authHeader = undefined;
    if (doAuth && window.nostr) {
        const authEvent = {
            id: null,
            pubkey: null,
            created_at: Math.floor(Date.now() / 1000),
            kind: 27235,
            tags: [['u', uploadUrl],['method', uploadMethod]],
            content: '',
            sig: null,
        };
        const signedAuthEvent = await window.nostr.signEvent(authEvent);
        let jsonAuthEvent = JSON.stringify(signedAuthEvent);
        let base64AuthEvent = btoa(jsonAuthEvent);
        authHeader = `Nostr: ${base64AuthEvent}`;
    }
    const headers = (authHeader ? {'Authorization':authHeader} : {});
    for (let file of files) { 
        const formData = new FormData(); 
        formData.append('file', file); 
        try { 
            //docs: https://github.com/nostrbuild/nostr.build/blob/main/api/v2/routes_upload.php
            const response = await fetch(
                uploadUrl, 
                { method: uploadMethod, body: formData, headers: headers}
            );
            const result = await response.json(); 
            if (result.status === 'success') { 
                urls.push(result.data[0].url);
                filesUploaded += 1;
            } 
        } catch (error) { 
            console.log("An error occurred during file upload", error);
        } 
    }
    return urls;

# Nostr Event format for options

{
    "id": "",
    "pubkey": "",
    "created_at": timestamp,
    "kind": 3939,
    "content": "",
    "tags": [
        ["L", "vote.pepe"],
        ["l", "voteoption"],
        ["t", "the category to use -- this value comes from config category"],
        ["r", "url pointing to the image to be voted on -- this value changes per url"],
    ],
    "sig": ""
}

# Nostr Event format for voting

{
    "id": "",
    "pubkey": "",
    "created_at": timestamp,
    "kind": 1212,
    "content": "",
    "tags": [
        ["L", "vote.pepe"],
        ["l", "vote"],
        ["t", "the category to use -- this value comes from config category"],
        ["e", "eventid of selected choice", "event id of unselected"],
    ],
    "sig": ""

}
