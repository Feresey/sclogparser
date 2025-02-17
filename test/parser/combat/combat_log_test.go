package parser_test

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Feresey/luxpanel/internal/parser/combat"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/connect.txt
	connectRaw string
	//go:embed testdata/start.txt
	startRaw string
	//go:embed testdata/damage.txt
	damageRaw string
	//go:embed testdata/heal.txt
	healRaw string
	//go:embed testdata/kill.txt
	killRaw string
	//go:embed testdata/finished.txt
	finishedRaw string
	//go:embed testdata/reward.txt
	rewardRaw string
	//go:embed testdata/apply_aura.txt
	applyAuraRaw string
)

func TestCombatConnectUnmarshal(t *testing.T) {
	now := time.Date(2023, time.January, 0, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name      string
		raw       string
		want      combat.ConnectToGameSession
		wantError bool
	}{
		{
			name: "ok",
			raw:  "19:32:58.666  CMBT   | ======= Connect to game session 50419619 =======",
			want: combat.ConnectToGameSession{
				LogTime:   time.Date(2023, 1, 0, 19, 32, 58, 666000000, time.Local),
				SessionID: 50419619,
			},
		},
		{
			name:      "cutted",
			raw:       "19:32:58.666  CMBT   | ======= Connect to game session 50419619",
			wantError: true,
		},
		{
			name:      "empty",
			raw:       "",
			wantError: true,
		},
		{
			name:      "wrong time",
			raw:       "19:32:58.66  CMBT   | ======= Connect to game session 50419619 =======",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			var val combat.ConnectToGameSession
			err := val.Unmarshal(tt.raw, now)
			if tt.wantError {
				r.Error(err)
				return
			} else {
				r.NoError(err)
			}
			r.Equal(tt.want, val)
		})
	}

	t.Run("raw", func(t *testing.T) {
		r := require.New(t)
		lines := strings.Split(connectRaw, "\n")

		for _, line := range lines {
			if line == "" {
				return
			}
			var val combat.ConnectToGameSession
			err := val.Unmarshal(line, now)
			r.NoError(err, line)
		}
	})
}

func New[T any](val T) *T {
	return &val
}

func TestCombatStartGameplayUnmarshal(t *testing.T) {
	now := time.Date(2023, time.January, 0, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name      string
		raw       string
		want      combat.StartGameplay
		wantError bool
	}{
		{
			name: "pve",
			raw:  `19:42:14.670  CMBT   | ======= Start PVE mission 'pve_raid_waterharvest_t5' map 'pve_raid_waterharvest' =======`,
			want: combat.StartGameplay{
				LogTime:      time.Date(2023, 1, 0, 19, 42, 14, 670000000, time.Local),
				GameMode:     "pve_raid_waterharvest_t5",
				MapName:      "pve_raid_waterharvest",
				ClientTeamID: nil,
			},
		},
		{
			name: "pvp",
			raw:  `20:21:02.744  CMBT   | ======= Start gameplay 'CaptureTheBase' map 's1420_ceres3_asteroidcity', local client team 1 =======`,
			want: combat.StartGameplay{
				LogTime:      time.Date(2023, 1, 0, 20, 21, 0o2, 744000000, time.Local),
				GameMode:     "CaptureTheBase",
				MapName:      "s1420_ceres3_asteroidcity",
				ClientTeamID: New(1),
			},
		},
		{
			name:      "cutted",
			raw:       "19:42:14.670  CMBT   | ======= Start PVE mission 'pve_raid_waterharvest_t5' map 'pve_raid_waterharvest'",
			wantError: true,
		},
		{
			name:      "empty",
			raw:       "",
			wantError: true,
		},
		{
			name:      "wrong time",
			raw:       `19:42:14.60  CMBT   | ======= Start PVE mission 'pve_raid_waterharvest_t5' map 'pve_raid_waterharvest' =======`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			var val combat.StartGameplay
			err := val.Unmarshal(tt.raw, now)
			if tt.wantError {
				r.Error(err)
				return
			} else {
				r.NoError(err)
			}
			r.Equal(tt.want, val)
		})
	}

	t.Run("raw", func(t *testing.T) {
		r := require.New(t)
		lines := strings.Split(startRaw, "\n")

		for _, line := range lines {
			if line == "" {
				return
			}
			var val combat.StartGameplay
			err := val.Unmarshal(line, now)
			r.NoError(err, line)
		}
	})
}

func TestCombatDamageUnmarshal(t *testing.T) {
	now := time.Date(2023, time.January, 0, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name      string
		raw       string
		want      combat.Damage
		wantError bool
	}{
		{
			name: "ok",
			raw:  `21:17:13.938  CMBT   | Damage        Gladiator|0000003117 ->           YanFei|0000167786  73.78 (h:0.00 s:73.78) Weapon_PlasmaBursts_T5_Rel EMP`,
			want: combat.Damage{
				LogTime:      time.Date(2023, 1, 0, 21, 17, 13, 938000000, time.Local),
				Initiator:    "Gladiator",
				InitiatorID:  3117,
				Recipient:    "YanFei",
				RecipientID:  167786,
				DamageTotal:  73.78,
				DamageHull:   0,
				DamageShield: 73.78,
				ActionSource: "Weapon_PlasmaBursts_T5_Rel",
				DamageModifiers: combat.DamageModifiersMap{
					"EMP": {},
				},
			},
		},
		{
			name: "no action",
			raw:  `19:44:04.074  CMBT   | Damage Megabomb_RW_BlackHole|0000000155 ->            tuman|0000000824   0.00 (h:0.00 s:0.00)  KINETIC`,
			want: combat.Damage{
				LogTime:      time.Date(2023, 1, 0, 19, 44, 4, 74000000, time.Local),
				Initiator:    "Megabomb_RW_BlackHole",
				InitiatorID:  155,
				Recipient:    "tuman",
				RecipientID:  824,
				DamageTotal:  0,
				DamageHull:   0,
				DamageShield: 0,
				ActionSource: "",
				DamageModifiers: combat.DamageModifiersMap{
					"KINETIC": {},
				},
			},
		},
		{
			name: "crash",
			raw:  `19:42:53.450  CMBT   | Damage            Py6Jl|0000000395 ->            Py6Jl|0000000395   0.00 (h:0.00 s:0.00) Weapon_OrbGun_T5_Epic EMP|PRIMARY_WEAPON|EXPLOSION <FriendlyFire>`,
			want: combat.Damage{
				LogTime:      time.Date(2023, 1, 0, 19, 42, 53, 450000000, time.Local),
				Initiator:    "Py6Jl",
				InitiatorID:  395,
				Recipient:    "Py6Jl",
				RecipientID:  395,
				DamageTotal:  0,
				DamageHull:   0,
				DamageShield: 0,
				ActionSource: "Weapon_OrbGun_T5_Epic",
				FriendlyFire: true,
				DamageModifiers: combat.DamageModifiersMap{
					"EMP": {}, "PRIMARY_WEAPON": {}, "EXPLOSION": {},
				},
			},
		},

		{
			name:      "cutted",
			raw:       `21:17:13.938  CMBT   | Damage        Gladiator|0000003117 ->           YanFei|0000167786  73.78 (h:0.00 s:73.78) Weapon_PlasmaBursts_T5_Rel`,
			wantError: true,
		},
		{
			name:      "empty",
			raw:       "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			var val combat.Damage
			err := val.Unmarshal(tt.raw, now)
			if tt.wantError {
				r.Error(err)
				return
			} else {
				r.NoError(err)
			}

			r.Equal(tt.want, val)
		})
	}

	t.Run("raw", func(t *testing.T) {
		r := require.New(t)
		lines := strings.Split(damageRaw, "\n")

		modifiers := make(map[combat.DamageModifier]int)
		for _, line := range lines {
			if line == "" {
				break
			}
			var val combat.Damage
			err := val.Unmarshal(line, now)
			r.NoError(err)

			for m := range val.DamageModifiers {
				modifiers[m]++
			}
		}

		fmt.Printf("%+#v\n", modifiers)
	})
}

func TestCombatHealUnmarshal(t *testing.T) {
	now := time.Date(2023, time.January, 0, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name      string
		raw       string
		want      combat.Heal
		wantError bool
	}{
		{
			name: "ok",
			raw:  `19:33:24.732  CMBT   | Heal            Feresey|0000000204 ->          Feresey|0000000204 244.00 Module_Lynx2Shield_T4_Epic`,
			want: combat.Heal{
				LogTime:      time.Date(2023, 1, 0, 19, 33, 24, 732000000, time.Local),
				Initiator:    "Feresey",
				InitiatorID:  204,
				Recipient:    "Feresey",
				RecipientID:  204,
				Heal:         244,
				ActionSource: "Module_Lynx2Shield_T4_Epic",
			},
		},

		{
			name:      "cutted",
			raw:       `19:33:24.732  CMBT   | Heal            Feresey|0000000204 ->          Feresey|0000000204 244.00`,
			wantError: true,
		},
		{
			name:      "empty",
			raw:       "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			var val combat.Heal
			err := val.Unmarshal(tt.raw, now)
			if tt.wantError {
				r.Error(err)
				return
			} else {
				r.NoError(err)
			}

			r.Equal(tt.want, val)
		})
	}

	t.Run("raw", func(t *testing.T) {
		r := require.New(t)
		lines := strings.Split(healRaw, "\n")

		for _, line := range lines {
			if line == "" {
				return
			}
			var val combat.Heal
			err := val.Unmarshal(line, now)
			r.NoError(err)
		}
	})
}

func TestCombatKillUnmarshal(t *testing.T) {
	now := time.Date(2023, time.January, 0, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name      string
		raw       string
		want      combat.Kill
		wantError bool
	}{
		{
			name: "player",
			raw:  `19:33:59.527  CMBT   | Killed Py6Jl      Ship_Race3_M_T2_Pirate|0000000248;      killer Feresey|0000000204 Weapon_Plasmagun_Heavy_T5_Pirate`,
			want: combat.Kill{
				LogTime:       time.Date(2023, 1, 0, 19, 33, 59, 527000000, time.Local),
				Initiator:     "Feresey",
				InitiatorID:   204,
				RecipientName: "Py6Jl",
				RecipientShip: "Ship_Race3_M_T2_Pirate",
				RecipientID:   248,
				ActionSource:  "Weapon_Plasmagun_Heavy_T5_Pirate",
			},
		},
		{
			name: "not player",
			raw:  `19:43:01.146  CMBT   | Killed Alien_Destroyer_Life_02_T5|0000001154;     killer Feresey|0000000766 Weapon_PlasmaWebLaser_T5_Epic`,
			want: combat.Kill{
				LogTime:         time.Date(2023, 1, 0, 19, 43, 1, 146000000, time.Local),
				Initiator:       "Feresey",
				InitiatorID:     766,
				RecipientObject: "Alien_Destroyer_Life_02_T5",
				RecipientID:     1154,
				ActionSource:    "Weapon_PlasmaWebLaser_T5_Epic",
			},
		},
		{
			name: "friendly fire",
			raw:  `19:46:16.971  CMBT   | Killed HealBot_Armor(Therm0Nuclear)|0000039068;	 killer Therm0Nuclear|0000039068 (suicide) <FriendlyFire>`,
			want: combat.Kill{
				LogTime:              time.Date(2023, 1, 0, 19, 46, 16, 971000000, time.Local),
				Initiator:            "Therm0Nuclear",
				InitiatorID:          39068,
				RecipientObjectOwner: "Therm0Nuclear",
				RecipientObjectName:  "HealBot_Armor",
				RecipientID:          39068,
				ActionSource:         "(suicide)",
				FriendlyFire:         true,
			},
		},
		{
			name: "killed object",
			raw:  `15:55:08.879  CMBT   | Killed SwarmPack3(MADEinHEAVEN)|0000056377;	 killer MADEinHEAVEN|0000056377 (suicide) <FriendlyFire>`,
			want: combat.Kill{
				LogTime:              time.Date(2023, 1, 0, 15, 55, 8, 879000000, time.Local),
				Initiator:            "MADEinHEAVEN",
				InitiatorID:          56377,
				RecipientObjectOwner: "MADEinHEAVEN",
				RecipientObjectName:  "SwarmPack3",
				RecipientID:          56377,
				ActionSource:         "(suicide)",
				FriendlyFire:         true,
			},
		},
		{
			name:      "cutted",
			raw:       `19:33:59.527  CMBT   | Killed Py6Jl      Ship_Race3_M_T2_Pirate|0000000248;      killer Feresey|`,
			wantError: true,
		},
		{
			name:      "empty",
			raw:       "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			var val combat.Kill
			err := val.Unmarshal(tt.raw, now)
			if tt.wantError {
				r.Error(err)
				return
			} else {
				r.NoError(err)
			}
			r.Equal(tt.want, val)
		})
	}

	t.Run("raw", func(t *testing.T) {
		r := require.New(t)
		lines := strings.Split(killRaw, "\n")

		for _, line := range lines {
			if line == "" {
				return
			}
			var val combat.Kill
			err := val.Unmarshal(line, now)
			r.NoError(err, line)
		}
	})
}

func TestCombatGameFinishedUnmarshal(t *testing.T) {
	now := time.Date(2023, time.January, 0, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name      string
		raw       string
		want      combat.FinishedGameplay
		wantError bool
	}{
		{
			name: "ok",
			raw:  `19:47:09.448  CMBT   | Gameplay finished. Winner team: 1(PVE_MISSION_COMPLETE_ALT_2). Finish reason: 'Mission complete'. Actual game time 275.9 sec`,
			want: combat.FinishedGameplay{
				LogTime:      time.Date(2023, 1, 0, 19, 47, 9, 448000000, time.Local),
				WinnerTeamID: 1,
				WinReason:    "PVE_MISSION_COMPLETE_ALT_2",
				FinishReason: "Mission complete",
				GameTime:     (275 * time.Second) + 900*time.Millisecond,
			},
		},
		{
			name:      "cutted",
			raw:       `19:47:09.448  CMBT   | Gameplay finished. Winner team: 1(PVE_MISSION_COMPLETE_ALT_2). Finish reason: 'Mission complete'. Actual game time`,
			wantError: true,
		},
		{
			name:      "empty",
			raw:       "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			var val combat.FinishedGameplay
			err := val.Unmarshal(tt.raw, now)
			if tt.wantError {
				r.Error(err)
				return
			} else {
				r.NoError(err)
			}
			r.Equal(tt.want, val)
		})
	}

	t.Run("raw", func(t *testing.T) {
		r := require.New(t)
		lines := strings.Split(finishedRaw, "\n")

		for _, line := range lines {
			if line == "" {
				return
			}
			var val combat.FinishedGameplay
			err := val.Unmarshal(line, now)
			r.NoError(err, line)
		}
	})
}

func TestCombatRewardUnmarshal(t *testing.T) {
	now := time.Date(2023, time.January, 0, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name      string
		raw       string
		want      combat.Reward
		wantError bool
	}{
		{
			name: "heal",
			raw:  `19:42:59.796  CMBT   | Reward idanceandkillyou Ship_Race3_L_T5_Faction1        	   1 experience                for heal Feresey`,
			want: combat.Reward{
				LogTime:    time.Date(2023, 1, 0, 19, 42, 59, 796000000, time.Local),
				Recipient:  "idanceandkillyou",
				Ship:       "Ship_Race3_L_T5_Faction1",
				Reward:     1,
				RewardType: "experience",
				Reason:     "heal Feresey",
			},
		},
		{
			name: "kill",
			raw:  `21:51:39.757  CMBT   | Reward        Gladiator Ship_Race5_H_OVERSEER_Rank15_13 	  46 effective points          for kill Khushal64n6`,
			want: combat.Reward{
				LogTime:    time.Date(2023, 1, 0, 21, 51, 39, 757000000, time.Local),
				Recipient:  "Gladiator",
				Ship:       "Ship_Race5_H_OVERSEER_Rank15_13",
				Reward:     46,
				RewardType: "effective points",
				Reason:     "kill Khushal64n6",
			},
		},
		{
			name: "victory",
			raw:  `21:51:41.115  CMBT   | Reward          Jigolee	722906 credits                   for victory`,
			want: combat.Reward{
				LogTime:    time.Date(2023, 1, 0, 21, 51, 41, 115000000, time.Local),
				Recipient:  "Jigolee",
				Ship:       "",
				Reward:     722906,
				RewardType: "credits",
				Reason:     "victory",
			},
		},

		{
			name: "strange kill",
			raw:  `21:51:31.180  CMBT   | Reward    PomogitePsixy Ship_Race2_M_T5_CraftUniq       	   9 effective points          for damage assist to kill `,
			want: combat.Reward{
				LogTime:    time.Date(2023, 1, 0, 21, 51, 31, 180000000, time.Local),
				Recipient:  "PomogitePsixy",
				Ship:       "Ship_Race2_M_T5_CraftUniq",
				Reward:     9,
				RewardType: "effective points",
				Reason:     "damage assist to kill ",
			},
		},
		{
			name: "cutted kill",
			raw:  `21:51:36.614  CMBT   | Reward      Khushal64n6 Ship_Race3_H_T5_Uniq            	 123 experience                for kill `,
			want: combat.Reward{
				LogTime:    time.Date(2023, 1, 0, 21, 51, 36, 614000000, time.Local),
				Recipient:  "Khushal64n6",
				Ship:       "Ship_Race3_H_T5_Uniq",
				Reward:     123,
				RewardType: "experience",
				Reason:     "kill ",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			var val combat.Reward
			err := val.Unmarshal(tt.raw, now)
			if tt.wantError {
				r.Error(err, "line: %s", tt.raw)
				return
			} else {
				r.NoError(err, "line: %s", tt.raw)
			}
			r.Equal(tt.want, val)
		})
	}

	t.Run("raw", func(t *testing.T) {
		r := require.New(t)
		lines := strings.Split(rewardRaw, "\n")

		for _, line := range lines {
			if line == "" {
				return
			}
			var val combat.Reward
			err := val.Unmarshal(line, now)
			r.NoError(err, line)
		}
	})
}

func TestCombatApplyAuraUnmarshal(t *testing.T) {
	now := time.Date(2023, time.January, 0, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name      string
		raw       string
		want      combat.ApplyAura
		wantError bool
	}{
		{
			name: "apply aura",
			raw:  `01:05:47.750  CMBT   | Apply aura 'TestKPMNormalizer_1' id 97 type AURA_RESIST_ALL to 'Taurefinne'`,
			want: combat.ApplyAura{
				LogTime:   time.Date(2023, 1, 0, 1, 5, 47, 750000000, time.Local),
				AuraName:  "TestKPMNormalizer_1",
				AuraID:    97,
				AuraType:  "AURA_RESIST_ALL",
				Recipient: "Taurefinne",
			},
		},
		{
			name: "apply aura to object",
			raw:  `23:04:41.679  CMBT   | Apply aura 'Spell_AdvancedHeal_T5_Epic' id 433 type AURA_HEALING_MOD to 'DestrMunition_SectorShield_T5_Mk3(Khushal64n6)'`,
			want: combat.ApplyAura{
				LogTime:     time.Date(2023, 1, 0, 23, 04, 41, 679000000, time.Local),
				AuraName:    "Spell_AdvancedHeal_T5_Epic",
				AuraID:      433,
				AuraType:    "AURA_HEALING_MOD",
				ObjectName:  "DestrMunition_SectorShield_T5_Mk3",
				ObjectOwner: "Khushal64n6",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			var val combat.ApplyAura
			err := val.Unmarshal(tt.raw, now)
			if tt.wantError {
				r.Error(err, "line: %s", tt.raw)
				return
			} else {
				r.NoError(err, "line: %s", tt.raw)
			}
			r.Equal(tt.want, val)
		})
	}

	t.Run("raw", func(t *testing.T) {
		r := require.New(t)
		lines := strings.Split(applyAuraRaw, "\n")

		for _, line := range lines {
			if line == "" {
				return
			}
			var val combat.ApplyAura
			err := val.Unmarshal(line, now)
			r.NoError(err, line)
		}
	})
}
