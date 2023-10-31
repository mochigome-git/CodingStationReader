// pkg/thd/text_insert.go
package thd

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"unicode"

	_ "github.com/lib/pq"
)

// Not used
// sub splits the input string into substrings based on non-letter and non-number characters.
func sub(contents string) []string {
	// Define a function 'f' that returns 'true' for non-letter and non-number runes.
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	// Use 'strings.FieldsFunc' to split the 'contents' string using the 'f' function as a separator.
	arr := strings.FieldsFunc(contents, f)
	fmt.Println(arr)
	// 'arr' is now a slice of substrings separated by non-letter and non-number characters.
	return arr
}

// Not used
// sub splits the input string into substrings based on non-letter and non-number characters,
// while preserving spaces, hyphens, and removing semicolons.
func sub2(contents string) []string {
	var result []string
	var currentSubStr string

	for _, char := range contents {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != ' ' && char != '-' {
			if currentSubStr != "" {
				// Replace multiple spaces with a single space within the current substring.
				currentSubStr = strings.Join(strings.Fields(currentSubStr), " ")
				result = append(result, strings.ReplaceAll(currentSubStr, ";", ""))
				currentSubStr = ""
			}
		}
		currentSubStr += string(char)
	}

	if currentSubStr != "" {
		// Replace multiple spaces with a single space within the last substring.
		currentSubStr = strings.Join(strings.Fields(currentSubStr), " ")
		result = append(result, strings.ReplaceAll(currentSubStr, ";", ""))
	}

	fmt.Println(result)
	// 'result' is now a slice of substrings separated by non-letter, non-number characters,
	// while preserving spaces and hyphens and removing semicolons.
	return result
}

// Not used
// sub splits the input string into substrings based on non-letter and non-number characters,
// while preserving hyphens.
func sub3(contents string) []string {
	// Define a function 'f' that returns 'true' for non-letter and non-number runes.
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	// Split the 'contents' string using the 'f' function as a separator, while preserving hyphens.
	arr := strings.FieldsFunc(contents, func(r rune) bool {
		return f(r) && r != '-'
	})
	fmt.Println(arr)
	// 'arr' is now a slice of substrings separated by non-letter and non-number characters,
	// while preserving hyphens.
	return arr
}

// Not used
// sub splits the input string into substrings based on non-letter and non-number characters,
// while preserving hyphens and the two big spaces.
func sub4(contents string) []string {
	// Define a function 'f' that returns 'true' for non-letter and non-number runes.
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	// Split the 'contents' string using a custom delimiter function.
	arr := strings.FieldsFunc(contents, func(r rune) bool {
		return f(r) && r != '-' && r != ' ' // Add ' ' as a delimiter
	})
	//fmt.Println(arr)
	// 'arr' is now a slice of substrings separated by non-letter and non-number characters,
	// while preserving hyphens and the two big spaces as separate strings.
	return arr
}

func Insert(host, port, user, password, dbname, pcName, joborder, contents string) (string, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("Failed to open database connection:", err)
		return "failed", err
	}
	defer db.Close()

	sqlStatement := `
	INSERT INTO coding_data (
		chip_tag,                            -- Slice[1]                   
		reserved,                            -- Slice[2]                   
		family_id,                           -- Slice[3]                   
		approved_hp_oem,                    -- Slice[4]                   
		oem_id,                              -- Slice[5]                   
		address_position,                    -- Slice[6]                   
		template_version_major,              -- Slice[7]                   
		template_version_minor,              -- Slice[8]                   
		tag_encryption_mask,                 -- Slice[9]                   
		person_revision,                     -- Slice[10]                  
		reserved_for_perso,                  -- Slice[11]                  
		ud0_fuse,                            -- Slice[12]                  
		dry_cartridge_sn_manufacture_site_id, -- Slice[13]                  
		dry_cartridge_sn_manufacture_line,    -- Slice[14]                  
		dry_cartridge_sn_manufacture_year,    -- Slice[15]                  
		dry_cartridge_sn_week_of_year,        -- Slice[16]                  
		dry_cartridge_sn_day_of_week,         -- Slice[17]                  
		dry_cartridge_sn_hour_of_day,         -- Slice[18]                  
		dry_cartridge_sn_minute_of_hour,      -- Slice[19]                  
		dry_cartridge_sn_sec_of_minute,        -- Slice[20]                  
		dry_cartridge_sn_process_position,     -- Slice[21]                  
		max_usable_cartridge_volume,          -- Slice[22]                  
		printer_lockdown_parition,           -- Slice[23]                  
		thermal_sense_resistor_tsr,          -- Slice[24]                  
		tsr_thermal_coeffcient_tcr,          -- Slice[25]                  
		bulk,                                -- Slice[26]                  
		ud1_fuse,                            -- Slice[27]                  
		cartridge_fill_sn_site_id,            -- Slice[28]                  
		cartridge_fill_sn_line,               -- Slice[29]                  
		cartridge_fill_sn_year,               -- Slice[30]                  
		cartridge_fill_sn_week_of_year,       -- Slice[31]                  
		cartridge_fill_sn_day_of_week,        -- Slice[32]                  
		cartridge_fill_sn_hour_of_day,        -- Slice[33]                  
		cartridge_fill_sn_minute_of_hour,     -- Slice[34]                  
		cartridge_fill_sn_sec_of_minute,       -- Slice[35]                  
		cartridge_fill_sn_process_position,    -- Slice[36]                  
		ink_formulator_id,                    -- Slice[37]                  
		ink_family,                           -- Slice[38]                  
		color_codes_general,                  -- Slice[39]                  
		color_codes_specific,                 -- Slice[40]                  
		ink_family_member,                    -- Slice[41]                  
		ink_id_number,                        -- Slice[42]                  
		ink_revision,                         -- Slice[43]                  
		ink_density,                          -- Slice[44]                  
		cartridge_distinction,                -- Slice[45]                  
		supply_key_size_descriptor,           -- Slice[46]                  
		shelf_life_weeks,                    -- Slice[47]                  
		shelf_life_days,                     -- Slice[48]                  
		installed_life_weeks,                -- Slice[49]                  
		installed_life_days,                 -- Slice[50]                  
		usable_ink_weight,                   -- Slice[51]                  
		altered_supply_notification_level,    -- Slice[52]                  
		firing_frequency,                    -- Slice[53]                  
		pulse_width_tpw,                     -- Slice[54]                  
		firing_voltage,                      -- Slice[55]                  
		turn_on_energy_toe,                  -- Slice[56]                  
		pulse_warming_temperature,            -- Slice[57]                  
		maximum_temperature,                 -- Slice[58]                  
		drop_volume,                         -- Slice[59]                  
		write_protect_fuse,                  -- Slice[60]                  
		_1st_platform_id,                    -- Slice[61]                   
		_1st_platform_manf_year,              -- Slice[62]                   
		_1st_platform_manf_week_of_year,      -- Slice[63]                   
		_1st_platform_mfg_country,           -- Slice[64]                   
		_1st_platform_fw_revision_major,      -- Slice[65]                   
		_1st_platform_fw_revision_minor,      -- Slice[66]                   
		_1st_install_cartridge_count,        -- Slice[67]                   
		cartridge_1st_install_year,          -- Slice[68]                   
		cartridge_1st_install_week_of_year,  -- Slice[69]                   
		cartridge_1st_install_day_of_week,   -- Slice[70]                   
		ink_level_gauge_resolution,           -- Slice[71]                   
		ud3_fuse,                            -- Slice[72]                   
		none,                                -- Slice[73]                   
		oem_defined_field_1,                 -- Slice[74]                   
		oem_defined_field_2,                 -- Slice[75]                   
		trademark_string,                    -- Slice[76]                   
		ud4_fuse,                            -- Slice[77]                   
		out_of_ink_bit,                      -- Slice[78]                   
		ilg_bits_1_25,                       -- Slice[79]                   
		ilg_bits_26_50,                      -- Slice[80]                   
		ilg_bits_51_75,                      -- Slice[81]                   
		ilg_bits_76_100,                     -- Slice[82]                   
		tiug_bits_1_25,                      -- Slice[83]                   
		tiug_bits_26_50,                     -- Slice[84]                   
		tiug_bits_51_75,                     -- Slice[85]                   
		tiug_bits_76_100,                    -- Slice[86]                   
		first_failure_code,                  -- Slice[87]                   
		altered_supply,                      -- Slice[88]                   
		user_acknowledge_altered_supply,     -- Slice[89]                   
		user_acknowledge_expired_ink,        -- Slice[90]                   
		faulty_replace_imeediately,          -- Slice[91]                   
		oem_defined_rw_or_field_1,           -- Slice[92]                   
		oem_defined_rw_or_field_2,           -- Slice[93]                   
		cartridge_mru_year,                  -- Slice[94]                   
		cartridge_mru_week_of_year,          -- Slice[95]                   
		cartridge_mru_day_of_week,           -- Slice[96]                   
		mru_platform_id,                     -- Slice[97]                   
		mru_platform_mfg_year,               -- Slice[98]                   
		mru_platform_mfg_week_of_year,       -- Slice[99]                   
		mru_platform_mfg_country,            -- Slice[100]                  
		mru_platform_fw_revision_major,      -- Slice[101]                  
		mru_platform_fw_revision_minor,      -- Slice[102]                  
		cartridge_insertion_count,           -- Slice[103]                  
		stall_insertion_count,               -- Slice[104]                  
		last_failure_code,                   -- Slice[105]                  
		last_user_reported_status,           -- Slice[106]                  
		marketing_data_revision,             -- Slice[107]                  
		oem_defined_rw_field_1,              -- Slice[108]                  
		oem_defined_rw_field_2,              -- Slice[109]                  
		ud7_fuse,                            -- Slice[110]                  
		extended_oem_id,                    -- Slice[111]                  
		hp_oem_ink_designator,               -- Slice[112]                  
		regionalization_id,                  -- Slice[113]                  
		cartridge_reorder_pn,                -- Slice[114]                  
		ud8_fuse,                            -- Slice[115]                  
		pcname,                              -- Slice[116]                  
		joborder                             -- Slice[117]                       
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
		$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38,
		$39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56,
		$57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74,
		$75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92,
		$93, $94, $95, $96, $97, $98, $99, $100, $101, $102, $103, $104, $105, $106, $107, $108,
		$109, $110, $111, $112, $113, $114, $115, $116, $117
	)
	RETURNING uuid
	`
	var uuid string
	Slice := sub4(contents)
	if len(Slice) == 0 {
		Slice = []string{"0000"}
	}

	arguments := make([]interface{}, 0, 116)
	for _, val := range Slice {

		arguments = append(arguments, val)
		//fmt.Printf("Slice %d: %v\n", i+1, val)
	}

	//Remove useless index
	totalToRemove := len(arguments) - 115
	indexToRemove := 115 // Subtract 1 since slices are zero-indexed

	//Combine the useful index with PCName and Joborder
	arguments = append(arguments[:indexToRemove], arguments[indexToRemove+totalToRemove:]...)
	arguments = append(arguments, pcName, joborder)

	err = db.QueryRow(sqlStatement, arguments...).Scan(&uuid)
	if err != nil {
		return "failed", fmt.Errorf("failed to execute SQL query: %w", err)
	}

	fmt.Println("Inserted row", uuid)

	return "complete", nil
}
