package proto

// question category list
const (
	Wb_quest_category_blank_fill   = "blank"
	Wb_quest_category_blank_len    = "blank_l"
	Wb_quest_category_option_chose = "chose"
)

var wbQuestionCategory = []string{
	Wb_quest_category_blank_fill,
	Wb_quest_category_blank_len,
	Wb_quest_category_option_chose,
}

func WbQuestionCagegories() []string {
	return wbQuestionCategory
}

// pvp mode type list
const (
	Wb_pvp_mode_normal = "normal"
	Wb_pvp_mode_race   = "race"
)

var wbPvpModes = []string{
	Wb_pvp_mode_normal,
	Wb_pvp_mode_race,
}

func WbPvpModes() []string {
	return wbPvpModes
}

// pvp difficulty type list
const (
	Wb_difficulty_normal  = "normal"
	Wb_difficulty_middle  = "middle"
	Wb_difficulty_extreme = "extreme"
)

var wbPvpDifficulties = []string{
	Wb_difficulty_normal,
	Wb_difficulty_middle,
	Wb_difficulty_extreme,
}

func WbPvpDifficulties() []string {
	return wbPvpDifficulties
}
