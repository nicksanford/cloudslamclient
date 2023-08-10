# Usage



```bash
export BASEURL=http://localhost:8080

# viam --base-url=$BASEURL organizations list
export ORGANIZATION_ID=

# viam --base-url=$BASEURL locations list
export LOCATION_ID=

# viam --base-url=$BASEURL robots list
export ROBOT_ID=

# viam --base-url=$BASEURL robot status --location $LOCATOIN_ID --organization $ORGANIZATION_ID --robot $ROBOT_ID
export ROBOT_PART_ID=

# go to app.viam.com > robot > setup > copy viam server config > secret
export ROBOT_PART_SECRET=

go run main.go -app_address http://localhost:8080  -loc_id <LOCATION_ID> -org_id <ORGANIZATION_ID> -robot_id <ROBOT_ID> -robot_part_id <ROBOT_PART_ID> -robot_part_secret <ROBOT_PART_SECRET>
```
