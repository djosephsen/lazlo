package lib

type ApiResponse struct {
	Bots          []Bot     `json:"bots,omitempty"`
	CacheVersion  string    `json:"cache_version,omitempty"`
	Channels      []Channel `json:"channels,omitempty"`
	Channel       Channel   `json:"channel,omitempty"`
	Groups        []Group   `json:"groups,omitempty"`
	Group         Group     `json:"group,omitempty"`
	IMs           []IM      `json:"ims,omitempty"`
	LatestEventTs string    `json:"latest_event_ts,omitempty"`
	Latest        string    `json:"latest,omitempty"`
	Ok            bool      `json:"ok,omitempty"`
	ReplyTo       int32     `json:"reply_to,omitempty"`
	Error         string    `json:"error,omitempty"`
	HasMore       bool      `json:"has_more,omitempty"`
	Self          Self      `json:"self,omitempty"`
	Team          Team      `json:"team,omitempty"`
	URL           string    `json:"url,omitempty"`
	Users         []User    `json:"users,omitempty"`
	User          User      `json:"user,omitempty"`
	Messages      []Event   `json:"messages,omitempty"`
}

type Event struct {
	ID           int32        `json:"id,omitempty"`
	Type         string       `json:"type,omitempty"`
	Channel      string       `json:"channel,omitempty"`
	Text         string       `json:"text,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
	User         string       `json:"user,omitempty"`
	UserName     string       `json:"username,omitempty"`
	BotID        string       `json:"bot_id,omitempty"`
	Subtype      string       `json:"subtype,omitempty"`
	Ts           string       `json:"ts,omitempty"`
	Broker       *Broker
	CallBackCode string `json:"callbackcode,omitempty"`
	Extra        map[string]interface{}
}

type Attachment struct {
	Fallback   string            `json:"fallback"`
	Color      string            `json:"color,omitempty"`
	Pretext    string            `json:"pretext,omitempty"`
	AuthorName string            `json:"author_name,omitempty"`
	AuthorLink string            `json:"author_link,omitempty"`
	AuthorIcon string            `json:"author_icon,omitempty"`
	Title      string            `json:"title,omitempty"`
	TitleLink  string            `json:"title_link,omitempty"`
	Text       string            `json:"text,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
	ImageUrl   string            `json:"image_url,omitempty"`
	ThumbUrl   string            `json:"thumb_url,omitempty"`
	MarkdownIn []string          `json:"mrkdwn_in,omitempty"`
}

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

type User struct {
	Color             string      `json:"color,omitempty"`
	Deleted           bool        `json:"deleted,omitempty"`
	HasFiles          bool        `json:"has_files,omitempty"`
	ID                string      `json:"id,omitempty"`
	IsAdmin           bool        `json:"is_admin,omitempty"`
	IsBot             bool        `json:"is_bot,omitempty"`
	IsOwner           bool        `json:"is_owner,omitempty"`
	IsPrimaryOwner    bool        `json:"is_primary_owner,omitempty"`
	IsRestricted      bool        `json:"is_restricted,omitempty"`
	IsUltraRestricted bool        `json:"is_ultra_restricted,omitempty"`
	Name              string      `json:"name,omitempty"`
	Phone             interface{} `json:"phone,omitempty"`
	Presence          string      `json:"presence,omitempty"`
	Profile           struct {
		Email              string `json:"email,omitempty"`
		FirstName          string `json:"first_name,omitempty"`
		Image192           string `json:"image_192,omitempty"`
		Image24            string `json:"image_24,omitempty"`
		Image32            string `json:"image_32,omitempty"`
		Image48            string `json:"image_48,omitempty"`
		Image72            string `json:"image_72,omitempty"`
		LastName           string `json:"last_name,omitempty"`
		Phone              string `json:"phone,omitempty"`
		RealName           string `json:"real_name,omitempty"`
		RealNameNormalized string `json:"real_name_normalized,omitempty"`
	} `json:"profile,omitempty"`
	RealName string      `json:"real_name,omitempty"`
	Skype    string      `json:"skype,omitempty"`
	Status   interface{} `json:"status,omitempty"`
	Tz       string      `json:"tz,omitempty"`
	TzLabel  string      `json:"tz_label,omitempty"`
	TzOffset float64     `json:"tz_offset,omitempty"`
	Extra    map[string]interface{}
}

type Channel struct {
	Created     float64  `json:"created,omitempty"`
	Creator     string   `json:"creator,omitempty"`
	ID          string   `json:"id,omitempty"`
	IsArchived  bool     `json:"is_archived,omitempty"`
	IsChannel   bool     `json:"is_channel,omitempty"`
	IsGeneral   bool     `json:"is_general,omitempty"`
	IsMember    bool     `json:"is_member,omitempty"`
	LastRead    string   `json:"last_read,omitempty"`
	Latest      Event    `json:"latest,omitempty"`
	Members     []string `json:"members,omitempty"`
	Name        string   `json:"name,omitempty"`
	Purpose     Topic    `json:"purpose,omitempty"`
	Topic       Topic    `json:"topic,omitempty"`
	UnreadCount float64  `json:"unread_count,omitempty"`
	Extra       map[string]interface{}
}

type Group struct {
	Created    int64    `json:"created,omitempty"`
	Creator    string   `json:"creator,omitempty"`
	ID         string   `json:"id,omitempty"`
	IsArchived bool     `json:"is_archived,omitempty"`
	IsGroup    bool     `json:"is_group,omitempty"`
	Members    []string `json:"members,omitempty"`
	Name       string   `json:"name,omitempty"`
	Purpose    Topic    `json:"purpose,omitempty"`
	Topic      Topic    `json:"topic,omitempty"`
	Extra      map[string]interface{}
}

type Topic struct {
	Creator string  `json:"creator,omitempty"`
	LastSet float64 `json:"last_set,omitempty"`
	Value   string  `json:"value,omitempty"`
}

type IM struct {
	Created       int64  `json:"created,omitempty"`
	ID            string `json:"id,omitempty"`
	IsIm          bool   `json:"is_im,omitempty"`
	IsUserDeleted bool   `json:"is_user_deleted,omitempty"`
	Latest        Event  `json:"latest,omitempty"`
	User          string `json:"user,omitempty"`
	Extra         map[string]interface{}
}

type Bot struct {
	Created       float64 `json:"created,omitempty"`
	Deleted       bool    `json:"deleted,omitempty"`
	Icons         Icon    `json:"icons,omitempty"`
	ID            string  `json:"id,omitempty"`
	IsIm          bool    `json:"is_im,omitempty"`
	IsUserDeleted bool    `json:"is_user_deleted,omitempty"`
	User          string  `json:"user,omitempty"`
	Name          string  `json:"name,omitempty"`
	Extra         map[string]interface{}
}

type Icon struct {
	Image192     string `json:"image_132,omitempty"`
	Image132     string `json:"image_132,omitempty"`
	Image102     string `json:"image_102,omitempty"`
	Image88      string `json:"image_88,omitempty"`
	Image68      string `json:"image_68,omitempty"`
	Image48      string `json:"image_48,omitempty"`
	Image44      string `json:"image_44,omitempty"`
	Image34      string `json:"image_34,omitempty"`
	ImageDefault bool   `json:"image_default,omitempty"`
}

type Self struct {
	Created        float64 `json:"created,omitempty"`
	ID             string  `json:"id,omitempty"`
	ManualPresence string  `json:"manual_presence,omitempty"`
	Name           string  `json:"name,omitempty"`
	Prefs          struct {
		AllChannelsLoud                 bool    `json:"all_channels_loud,omitempty"`
		ArrowHistory                    bool    `json:"arrow_history,omitempty"`
		AtChannelSuppressedChannels     string  `json:"at_channel_suppressed_channels,omitempty"`
		AutoplayChatSounds              bool    `json:"autoplay_chat_sounds,omitempty"`
		Collapsible                     bool    `json:"collapsible,omitempty"`
		CollapsibleByClick              bool    `json:"collapsible_by_click,omitempty"`
		ColorNamesInList                bool    `json:"color_names_in_list,omitempty"`
		CommaKeyPrefs                   bool    `json:"comma_key_prefs,omitempty"`
		ConvertEmoticons                bool    `json:"convert_emoticons,omitempty"`
		DisplayRealNamesOverride        float64 `json:"display_real_names_override,omitempty"`
		DropboxEnabled                  bool    `json:"dropbox_enabled,omitempty"`
		EmailAlerts                     string  `json:"email_alerts,omitempty"`
		EmailAlertsSleepUntil           float64 `json:"email_alerts_sleep_until,omitempty"`
		EmailMisc                       bool    `json:"email_misc,omitempty"`
		EmailWeekly                     bool    `json:"email_weekly,omitempty"`
		EmojiMode                       string  `json:"emoji_mode,omitempty"`
		EnterIsSpecialInTbt             bool    `json:"enter_is_special_in_tbt,omitempty"`
		ExpandInlineImgs                bool    `json:"expand_inline_imgs,omitempty"`
		ExpandInternalInlineImgs        bool    `json:"expand_internal_inline_imgs,omitempty"`
		ExpandNonMediaAttachments       bool    `json:"expand_non_media_attachments,omitempty"`
		ExpandSnippets                  bool    `json:"expand_snippets,omitempty"`
		FKeySearch                      bool    `json:"f_key_search,omitempty"`
		FullTextExtracts                bool    `json:"full_text_extracts,omitempty"`
		FuzzyMatching                   bool    `json:"fuzzy_matching,omitempty"`
		GraphicEmoticons                bool    `json:"graphic_emoticons,omitempty"`
		GrowlsEnabled                   bool    `json:"growls_enabled,omitempty"`
		HasCreatedChannel               bool    `json:"has_created_channel,omitempty"`
		HasInvited                      bool    `json:"has_invited,omitempty"`
		HasUploaded                     bool    `json:"has_uploaded,omitempty"`
		HighlightWords                  string  `json:"highlight_words,omitempty"`
		KKeyOmnibox                     bool    `json:"k_key_omnibox,omitempty"`
		LastSnippetType                 string  `json:"last_snippet_type,omitempty"`
		LoudChannels                    string  `json:"loud_channels,omitempty"`
		LoudChannelsSet                 string  `json:"loud_channels_set,omitempty"`
		LsDisabled                      bool    `json:"ls_disabled,omitempty"`
		MacSpeakSpeed                   float64 `json:"mac_speak_speed,omitempty"`
		MacSpeakVoice                   string  `json:"mac_speak_voice,omitempty"`
		MacSsbBounce                    string  `json:"mac_ssb_bounce,omitempty"`
		MacSsbBullet                    bool    `json:"mac_ssb_bullet,omitempty"`
		MarkMsgsReadImmediately         bool    `json:"mark_msgs_read_immediately,omitempty"`
		MessagesTheme                   string  `json:"messages_theme,omitempty"`
		MuteSounds                      bool    `json:"mute_sounds,omitempty"`
		MutedChannels                   string  `json:"muted_channels,omitempty"`
		NeverChannels                   string  `json:"never_channels,omitempty"`
		NewMsgSnd                       string  `json:"new_msg_snd,omitempty"`
		NoCreatedOverlays               bool    `json:"no_created_overlays,omitempty"`
		NoJoinedOverlays                bool    `json:"no_joined_overlays,omitempty"`
		NoMacssb1Banner                 bool    `json:"no_macssb1_banner,omitempty"`
		NoTextInNotifications           bool    `json:"no_text_in_notifications,omitempty"`
		ObeyInlineImgLimit              bool    `json:"obey_inline_img_limit,omitempty"`
		PagekeysHandled                 bool    `json:"pagekeys_handled,omitempty"`
		PostsFormattingGuide            bool    `json:"posts_formatting_guide,omitempty"`
		PrivacyPolicySeen               bool    `json:"privacy_policy_seen,omitempty"`
		PromptedForEmailDisabling       bool    `json:"prompted_for_email_disabling,omitempty"`
		PushAtChannelSuppressedChannels string  `json:"push_at_channel_suppressed_channels,omitempty"`
		PushDmAlert                     bool    `json:"push_dm_alert,omitempty"`
		PushEverything                  bool    `json:"push_everything,omitempty"`
		PushIdleWait                    float64 `json:"push_idle_wait,omitempty"`
		PushLoudChannels                string  `json:"push_loud_channels,omitempty"`
		PushLoudChannelsSet             string  `json:"push_loud_channels_set,omitempty"`
		PushMentionAlert                bool    `json:"push_mention_alert,omitempty"`
		PushMentionChannels             string  `json:"push_mention_channels,omitempty"`
		PushSound                       string  `json:"push_sound,omitempty"`
		RequireAt                       bool    `json:"require_at,omitempty"`
		SearchExcludeBots               bool    `json:"search_exclude_bots,omitempty"`
		SearchExcludeChannels           string  `json:"search_exclude_channels,omitempty"`
		SearchOnlyMyChannels            bool    `json:"search_only_my_channels,omitempty"`
		SearchSort                      string  `json:"search_sort,omitempty"`
		SeenChannelMenuTipCard          bool    `json:"seen_channel_menu_tip_card,omitempty"`
		SeenChannelsTipCard             bool    `json:"seen_channels_tip_card,omitempty"`
		SeenDomainInviteReminder        bool    `json:"seen_domain_invite_reminder,omitempty"`
		SeenFlexpaneTipCard             bool    `json:"seen_flexpane_tip_card,omitempty"`
		SeenMemberInviteReminder        bool    `json:"seen_member_invite_reminder,omitempty"`
		SeenMessageInputTipCard         bool    `json:"seen_message_input_tip_card,omitempty"`
		SeenSearchInputTipCard          bool    `json:"seen_search_input_tip_card,omitempty"`
		SeenSsbPrompt                   bool    `json:"seen_ssb_prompt,omitempty"`
		SeenTeamMenuTipCard             bool    `json:"seen_team_menu_tip_card,omitempty"`
		SeenUserMenuTipCard             bool    `json:"seen_user_menu_tip_card,omitempty"`
		SeenWelcome2                    bool    `json:"seen_welcome_2,omitempty"`
		ShowMemberPresence              bool    `json:"show_member_presence,omitempty"`
		ShowTyping                      bool    `json:"show_typing,omitempty"`
		SidebarBehavior                 string  `json:"sidebar_behavior,omitempty"`
		SidebarTheme                    string  `json:"sidebar_theme,omitempty"`
		SidebarThemeCustomValues        string  `json:"sidebar_theme_custom_values,omitempty"`
		SnippetEditorWrapLongLines      bool    `json:"snippet_editor_wrap_long_lines,omitempty"`
		SpeakGrowls                     bool    `json:"speak_growls,omitempty"`
		SsEmojis                        bool    `json:"ss_emojis,omitempty"`
		StartScrollAtOldest             bool    `json:"start_scroll_at_oldest,omitempty"`
		TabUiReturnSelects              bool    `json:"tab_ui_return_selects,omitempty"`
		Time24                          bool    `json:"time24,omitempty"`
		Tz                              string  `json:"tz,omitempty"`
		UserColors                      string  `json:"user_colors,omitempty"`
		WebappSpellcheck                bool    `json:"webapp_spellcheck,omitempty"`
		WelcomeMessageHidden            bool    `json:"welcome_message_hidden,omitempty"`
		WinSsbBullet                    bool    `json:"win_ssb_bullet,omitempty"`
	} `json:"prefs,omitempty"`
}

type Team struct {
	Domain            string  `json:"domain,omitempty"`
	EmailDomain       string  `json:"email_domain,omitempty"`
	Icon              Icon    `json:"icon,omitempty"`
	ID                string  `json:"id,omitempty"`
	MsgEditWindowMins float64 `json:"msg_edit_window_mins,omitempty"`
	Name              string  `json:"name,omitempty"`
	OverStorageLimit  bool    `json:"over_storage_limit,omitempty"`
	Prefs             struct {
		AllowMessageDeletion   bool     `json:"allow_message_deletion,omitempty"`
		DefaultChannels        []string `json:"default_channels,omitempty"`
		DisplayRealNames       bool     `json:"display_real_names,omitempty"`
		DmRetentionDuration    float64  `json:"dm_retention_duration,omitempty"`
		DmRetentionType        float64  `json:"dm_retention_type,omitempty"`
		GatewayAllowIrcPlain   float64  `json:"gateway_allow_irc_plain,omitempty"`
		GatewayAllowIrcSsl     float64  `json:"gateway_allow_irc_ssl,omitempty"`
		GatewayAllowXmppSsl    float64  `json:"gateway_allow_xmpp_ssl,omitempty"`
		GroupRetentionDuration float64  `json:"group_retention_duration,omitempty"`
		GroupRetentionType     float64  `json:"group_retention_type,omitempty"`
		HideReferers           bool     `json:"hide_referers,omitempty"`
		MsgEditWindowMins      float64  `json:"msg_edit_window_mins,omitempty"`
		RequireAtForMention    float64  `json:"require_at_for_mention,omitempty"`
		RetentionDuration      float64  `json:"retention_duration,omitempty"`
		RetentionType          float64  `json:"retention_type,omitempty"`
		WhoCanArchiveChannels  string   `json:"who_can_archive_channels,omitempty"`
		WhoCanAtChannel        string   `json:"who_can_at_channel,omitempty"`
		WhoCanAtEveryone       string   `json:"who_can_at_everyone,omitempty"`
		WhoCanCreateChannels   string   `json:"who_can_create_channels,omitempty"`
		WhoCanCreateGroups     string   `json:"who_can_create_groups,omitempty"`
		WhoCanKickChannels     string   `json:"who_can_kick_channels,omitempty"`
		WhoCanKickGroups       string   `json:"who_can_kick_groups,omitempty"`
		WhoCanPostGeneral      string   `json:"who_can_post_general,omitempty"`
	} `json:"prefs,omitempty"`
	Extra map[string]interface{}
}
