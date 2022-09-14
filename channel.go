package main

import (
	"fmt"
	"log"
	"text/tabwriter"

	"github.com/flatcar/updateservicectl/client/update/v1"
)

var (
	channelFlags struct {
		appId   StringFlag
		channel StringFlag
		version StringFlag
		publish bool
	}

	cmdChannel = &Command{
		Name:    "channel",
		Summary: "Manage channels for an application",
		Subcommands: []*Command{
			cmdChannelList,
			cmdChannelUpdate,
			cmdChannelCreate,
			cmdChannelDelete,
		},
	}

	cmdChannelList = &Command{
		Name:        "channel list",
		Usage:       "[OPTION]...",
		Description: `List all channels for an application.`,
		Run:         channelList,
	}

	cmdChannelCreate = &Command{
		Name:  "channel create",
		Usage: "[OPTION]...",
		Description: `Given an application ID (--app-id) and channel (--channel),
you can create a new channel.`,
		Run: channelCreate,
	}

	cmdChannelUpdate = &Command{
		Name:    "channel update",
		Usage:   "[OPTION]...",
		Summary: `Update the version and publish state for an application channel.`,
		Description: `Given an application ID (--app-id) and channel (--channel),
you can change the channel to a new version (--version), or set the publish state (--publish).`,
		Run: channelUpdate,
	}

	cmdChannelDelete = &Command{
		Name:        "channel delete",
		Usage:       "[OPTION]...",
		Summary:     `Delete an application channel.`,
		Description: `Deletes the channel with matching application ID (--app-id) and channel (--channel).`,
		Run:         channelDelete,
	}
)

func init() {
	cmdChannelList.Flags.Var(&channelFlags.appId, "app-id", "The application ID to list the channels of.")

	cmdChannelCreate.Flags.Var(&channelFlags.appId, "app-id", "The application ID that the channel belongs to.")
	cmdChannelCreate.Flags.Var(&channelFlags.channel, "channel", "The channel to create.")
	cmdChannelCreate.Flags.BoolVar(&channelFlags.publish, "publish", false, "Publish or unpublish the channel.")
	cmdChannelCreate.Flags.Var(&channelFlags.version, "version", "The version to create the channel to.")

	cmdChannelUpdate.Flags.Var(&channelFlags.appId, "app-id", "The application ID that the channel belongs to.")
	cmdChannelUpdate.Flags.Var(&channelFlags.channel, "channel", "The channel to update.")
	cmdChannelUpdate.Flags.BoolVar(&channelFlags.publish, "publish", false, "Publish or unpublish the channel.")
	cmdChannelUpdate.Flags.Var(&channelFlags.version, "version", "The version to update the channel to.")

	cmdChannelDelete.Flags.Var(&channelFlags.appId, "app-id", "The application ID that the channel belongs to.")
	cmdChannelDelete.Flags.Var(&channelFlags.channel, "channel", "The channel to update.")
}

const channelHeader = "Label\tVersion\tPublish\tUpstream\n"

func formatChannel(channel *update.AppChannel) string {
	return fmt.Sprintf("%s\t%s\t%t\t%s\n", channel.Label, channel.Version, channel.Publish, channel.Upstream)
}

func channelList(args []string, service *update.Service, out *tabwriter.Writer) int {
	if channelFlags.appId.Get() == nil {
		return ERROR_USAGE
	}

	listCall := service.Channel.List(channelFlags.appId.String())
	list, err := listCall.Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(out, channelHeader)
	for _, channel := range list.Items {
		fmt.Fprintf(out, "%s", formatChannel(channel))
	}
	out.Flush()
	return OK
}

func channelCreate(args []string, service *update.Service, out *tabwriter.Writer) int {
	if channelFlags.version.Get() == nil || channelFlags.appId.Get() == nil || channelFlags.channel.Get() == nil {
		return ERROR_USAGE
	}

	channelReq := &update.ChannelRequest{
		Version: *channelFlags.version.Get(),
		Publish: channelFlags.publish,
		Label:   *channelFlags.channel.Get(),
		AppId:   *channelFlags.appId.Get(),
	}

	call := service.Channel.Insert(*channelFlags.channel.Get(), channelReq)
	channel, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(out, channelHeader)
	fmt.Fprintf(out, "%s", formatChannel(channel))
	out.Flush()
	return OK
}

func channelUpdate(args []string, service *update.Service, out *tabwriter.Writer) int {
	if channelFlags.version.Get() == nil || channelFlags.appId.Get() == nil || channelFlags.channel.Get() == nil {
		return ERROR_USAGE
	}

	channelReq := &update.ChannelRequest{Version: *channelFlags.version.Get(), Publish: channelFlags.publish}

	call := service.Channel.Update(channelFlags.appId.String(), channelFlags.channel.String(), channelReq)
	channel, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(out, channelHeader)
	fmt.Fprintf(out, "%s", formatChannel(channel))
	out.Flush()
	return OK
}

func channelDelete(args []string, service *update.Service, out *tabwriter.Writer) int {
	if channelFlags.appId.Get() == nil || channelFlags.channel.Get() == nil {
		return ERROR_USAGE
	}

	call := service.Channel.Delete(channelFlags.appId.String(), channelFlags.channel.String())
	_, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(out, "deleted channel: %s, for application: %s\n", channelFlags.channel.String(), channelFlags.appId.String())
	out.Flush()
	return OK
}
