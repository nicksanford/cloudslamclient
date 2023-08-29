package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/edaniels/golog"
	"github.com/google/uuid"
	"go.uber.org/zap"
	pb "go.viam.com/api/app/cloudslam/v1"
	v1 "go.viam.com/api/app/cloudslam/v1"
	rutils "go.viam.com/rdk/utils"
	"go.viam.com/utils"
	"go.viam.com/utils/rpc"
)

type Args struct {
	AppAddress      string `flag:"app_address,required,usage=app_address"`
	RobotPartID     string `flag:"robot_part_id,required,usage=robot_part_id"`
	RobotID         string `flag:"robot_id,required,usage=robot_id"`
	RobotPartSecret string `flag:"robot_part_secret,required,usage=robot_part_secret"`
	OrganizationID  string `flag:"org_id,required,usage=organization_id"`
	LocationID      string `flag:"loc_id,required,usage=location_id"`
}

var logger = golog.NewDevelopmentLogger("cloudslamclient")

func main() {
	utils.ContextualMain(runMain, logger)
}

func runMain(ctx context.Context, rawArgs []string, logger *zap.SugaredLogger) error {
	var parsedArgs Args
	if err := utils.ParseFlags(rawArgs, &parsedArgs); err != nil {
		return err
	}

	conn, err := createNewGRPCClient(ctx, parsedArgs, logger)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewPackageClientFromConn(conn)
	u := uuid.New()

	reqStart := &v1.StartMappingSessionRequest{
		OrganizationId: parsedArgs.OrganizationID,
		LocationId:     parsedArgs.LocationID,
		RobotId:        parsedArgs.RobotID,
		MapName:        u.String(),
	}

	respStart, err := client.StartMappingSession(ctx, reqStart)
	if err != nil {
		return err
	}
	fmt.Printf("StartMappingSessionResponse: %#v\n", respStart.String())

	reqGetMetadata := &v1.GetMappingSessionMetadataByIDRequest{SessionId: respStart.SessionId}
	respGetMetadata, err := client.GetMappingSessionMetadataByID(ctx, reqGetMetadata)
	if err != nil {
		return err
	}
	fmt.Printf("GetMappingSessionMetadataByIDResponse: %#v\n", respGetMetadata.String())

	return nil
}

// createNewGRPCClient creates a new grpc cloud configured to communicate with the robot service based on the cloud config given.
func createNewGRPCClient(ctx context.Context, args Args, logger golog.Logger) (rpc.ClientConn, error) {
	u, err := url.Parse(args.AppAddress)
	if err != nil {
		return nil, err
	}

	dialOpts := make([]rpc.DialOption, 0, 2)
	dialOpts = append(dialOpts, rpc.WithEntityCredentials(args.RobotPartID,
		rpc.Credentials{
			Type:    rutils.CredentialsTypeRobotSecret,
			Payload: args.RobotPartSecret,
		},
	))

	if u.Scheme == "http" {
		dialOpts = append(dialOpts, rpc.WithInsecure())
	}

	return rpc.DialDirectGRPC(ctx, u.Host, logger, dialOpts...)
}

// NewPackageClientFromConn creates a new CloudSLAMClient.
func NewPackageClientFromConn(conn rpc.ClientConn) pb.CloudSLAMServiceClient {
	c := pb.NewCloudSLAMServiceClient(conn)
	return c
}
