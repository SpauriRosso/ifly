package music

import (
	"image"
	"reflect"
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// countLeadNotes takes a score and returns the number of melodic events used to compare visual densities.
func countLeadNotes(plan MusicPlan) int {
	count := 0
	for _, note := range plan.Notes {
		if note.Voice == "lead" {
			count++
		}
	}
	return count
}

// countVoiceInSection takes a score, voice, and beat range and returns the matching event count for arrangement checks.
func countVoiceInSection(plan MusicPlan, voice string, start, end float64) int {
	count := 0
	for _, note := range plan.Notes {
		if note.Voice == voice && note.Beat >= start && note.Beat < end {
			count++
		}
	}
	return count
}

// TestMusicServiceCompose checks that one artwork has a repeatable score and opposing scenes produce different scores.
func TestMusicServiceCompose(t *testing.T) {
	dataURL := encodeTestImage(t, image.NewRGBA(image.Rect(0, 0, 16, 16)))
	service := &MusicService{}
	calm := semantic.SemanticAnalysis{Description: "A serene landscape at sunset with a deer."}
	first, err := service.Compose(dataURL, calm)
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Compose(dataURL, calm)
	if err != nil || !reflect.DeepEqual(first, second) {
		t.Fatalf("same artwork changed score: %v", err)
	}
	if first.Mood != "calm" || first.TempoBPM < 55 || first.BeatsPerLoop != len(formEpisodes)*16 || len(first.Sections) != len(formEpisodes) || len(first.Notes) < 300 || len(first.Notes) > 1200 {
		t.Fatalf("invalid calm score: mood=%s tempo=%d beats=%d sections=%d notes=%d", first.Mood, first.TempoBPM, first.BeatsPerLoop, len(first.Sections), len(first.Notes))
	}
	threat, err := service.Compose(dataURL, semantic.SemanticAnalysis{Description: "A threatening violent scene."})
	if err != nil {
		t.Fatal(err)
	}
	if threat.Mood != "tense" || threat.Mode == first.Mode || threat.TempoBPM == first.TempoBPM || reflect.DeepEqual(threat.Notes, first.Notes) {
		t.Fatal("opposing scene produced the same score")
	}
	for index, note := range first.Notes {
		if note.Beat < 0 || note.Beat >= float64(first.BeatsPerLoop) || note.Duration <= 0 || note.MIDI < 36 || note.MIDI > 90 || note.Velocity <= 0 || note.Velocity > 1 {
			t.Fatalf("unsafe note: %+v", note)
		}
		if index > 0 && note.Beat < first.Notes[index-1].Beat {
			t.Fatal("score events are not sorted")
		}
		if note.Voice == "lead" && note.Beat > float64(first.BeatsPerLoop-4) {
			t.Fatal("the concluding bar should leave room for silence")
		}
	}
}

// TestPondSceneSelectsGentleCharacter checks that visible water lilies and reflective light select a calm palette without artist recognition.
func TestPondSceneSelectsGentleCharacter(t *testing.T) {
	character := selectMusicalCharacter(semantic.SemanticAnalysis{
		Description: "A painting of a bridge over a pond with water lilies and trees.",
		Atmosphere:  "soft reflective light over still water",
		Subjects:    []string{"bridge", "water lilies", "trees"},
	})
	if character.mood != "calm" {
		t.Fatalf("pond scene did not select a gentle palette: %s", character.mood)
	}
}

// TestLeadAndHumResolveWithinHarmony checks that melodic anchors and vocal notes fit each bar's chord and weak dissonances resolve by step.
func TestLeadAndHumResolveWithinHarmony(t *testing.T) {
	for _, description := range []string{"A pond with water lilies.", "A threatening storm.", "A lively dancing scene."} {
		interpretation := semantic.SemanticAnalysis{Description: description}
		character := selectMusicalCharacter(interpretation)
		for center := byte(0); center < 12; center++ {
			var digest [32]byte
			digest[0], digest[1], digest[2] = center, center%2, center%3
			plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, interpretation, digest)
			root := 48 + int(center)
			for bar := 0; bar < plan.BeatsPerLoop/4; bar++ {
				chord := harmonyForBar(character, bar, formEpisodes[bar/4].role, plan.BeatsPerLoop/4-1)
				members := chordPitchClasses(root+chord.rootOffset, chord.third)
				lead := []MusicNote{}
				for _, note := range plan.Notes {
					if note.Beat < float64(bar*4) || note.Beat >= float64((bar+1)*4) {
						continue
					}
					if note.Voice == "lead" {
						lead = append(lead, note)
					}
					if note.Voice == "hum" && !pitchBelongsToChord(note.MIDI, members) {
						t.Fatalf("unresolved hum in %q, bar %d: %+v", description, bar, note)
					}
				}
				for index, note := range lead {
					if pitchBelongsToChord(note.MIDI, members) {
						continue
					}
					if index == 0 || index+1 == len(lead) || absInt(note.MIDI-lead[index+1].MIDI) > 2 {
						t.Fatalf("unresolved lead in %q, center %d, bar %d: %+v", description, center, bar, lead)
					}
				}
			}
		}
	}
}

// absInt takes an integer and returns its nonnegative magnitude for interval checks.
func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// TestComposeMusicUsesVisualComplexity checks that added visual detail changes melodic density for the same interpreted scene.
func TestComposeMusicUsesVisualComplexity(t *testing.T) {
	interpretation := semantic.SemanticAnalysis{Description: "A serene landscape."}
	quiet := composeMusic(analysis.StructuralAnalysis{Complexity: 0.1}, interpretation, [32]byte{})
	busy := composeMusic(analysis.StructuralAnalysis{Complexity: 0.8}, interpretation, [32]byte{})
	if countLeadNotes(busy) <= countLeadNotes(quiet) || busy.TempoBPM <= quiet.TempoBPM {
		t.Fatal("visual complexity did not alter density and pace")
	}
}

// TestComposeMusicVariesTempo checks that the bridge accelerates and the conclusion settles for contrasting scenes.
func TestComposeMusicVariesTempo(t *testing.T) {
	for _, description := range []string{"A serene landscape.", "A threatening scene.", "A lively dancing scene."} {
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.6}, semantic.SemanticAnalysis{Description: description}, [32]byte{})
		if plan.Sections[0].TempoBPM != plan.TempoBPM || plan.Sections[2].TempoBPM <= plan.Sections[0].TempoBPM || plan.Sections[3].TempoBPM >= plan.Sections[2].TempoBPM || plan.Sections[4].TempoBPM <= plan.Sections[3].TempoBPM || plan.Sections[len(plan.Sections)-1].TempoBPM >= plan.Sections[9].TempoBPM {
			t.Fatalf("tempo arc did not develop for %q: %+v", description, plan.Sections)
		}
		for _, section := range plan.Sections {
			if section.TempoBPM < 40 || section.TempoBPM > 180 {
				t.Fatalf("tempo outside playable range: %+v", section)
			}
		}
	}
}

// TestLongFormDuration checks that calm and active scenes both produce a complete performance longer than two minutes.
func TestLongFormDuration(t *testing.T) {
	for _, description := range []string{"A pond with water lilies.", "A threatening storm.", "A lively dancing scene."} {
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.7}, semantic.SemanticAnalysis{Description: description}, [32]byte{})
		seconds := 0.0
		for _, section := range plan.Sections {
			seconds += (section.EndBeat - section.StartBeat) * 60 / float64(section.TempoBPM)
		}
		if seconds < 120 || len(plan.Sections) != 16 {
			t.Fatalf("short or incomplete form for %q: %.1f seconds, %d phrases", description, seconds, len(plan.Sections))
		}
	}
}

// TestComposeMusicResolvesHome checks that a tense phrase ends on its tonic instead of looping from an unresolved dominant.
func TestComposeMusicResolvesHome(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{}, semantic.SemanticAnalysis{Description: "A threatening scene."}, [32]byte{})
	for _, note := range plan.Notes {
		if note.Beat == float64(plan.BeatsPerLoop-4) && note.Voice == "pad" && note.MIDI%12 == 0 {
			return
		}
	}
	t.Fatal("final harmony did not resolve to the tonic")
}

// TestHarmonyForBarPreparesReturns checks that complete phrases move to a dominant before returning to the tonic.
func TestHarmonyForBarPreparesReturns(t *testing.T) {
	for name, character := range musicalStyles {
		for _, section := range []int{0, 4} {
			cadence := harmonyForBar(character, section*4+3, section, len(formEpisodes)*4-1)
			if cadence.rootOffset != character.bridge[3].rootOffset || cadence.rootOffset == 0 {
				t.Fatalf("%s section %d lacks a dominant preparation: %+v", name, section, cadence)
			}
		}
		finalBar := len(formEpisodes)*4 - 1
		final := harmonyForBar(character, finalBar, 5, finalBar)
		if final.rootOffset != 0 {
			t.Fatalf("%s ending does not resolve: %+v", name, final)
		}
	}
}

// TestVoiceChordPrefersNearInversion checks that chord voices move less than a root-position block change.
func TestVoiceChordPrefersNearInversion(t *testing.T) {
	first := voiceChord(48, 4, [3]int{}, false)
	second := voiceChord(53, 4, first, true)
	motion := 0
	for note := range first {
		delta := first[note] - second[note]
		if delta < 0 {
			delta = -delta
		}
		motion += delta
	}
	if motion >= 15 {
		t.Fatalf("chord motion is too large: %v to %v", first, second)
	}
}

// TestArpeggioBuildsAndRecedes checks that the counterline enters after the opening and withdraws before the final cadence.
func TestArpeggioBuildsAndRecedes(t *testing.T) {
	calm := composeMusic(analysis.StructuralAnalysis{}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	tense := composeMusic(analysis.StructuralAnalysis{}, semantic.SemanticAnalysis{Description: "A threatening scene."}, [32]byte{})
	if countVoiceInSection(calm, "arp", 0, 16) != 0 || countVoiceInSection(calm, "arp", 32, 48) == 0 || countVoiceInSection(calm, "arp", 48, 64) != 0 || countVoiceInSection(calm, "arp", float64(calm.BeatsPerLoop-8), float64(calm.BeatsPerLoop)) != 0 {
		t.Fatal("the calm counterline did not follow the phrase arc")
	}
	if countVoiceInSection(tense, "arp", 32, 48) <= countVoiceInSection(calm, "arp", 32, 48) {
		t.Fatal("the tense bridge should carry more counterline motion")
	}
}

// TestBassLineChangesItsResponse checks that an active bass voice uses different chord tones and the calm bridge adds movement.
func TestBassLineChangesItsResponse(t *testing.T) {
	character := musicalStyles["energetic"]
	intervals := map[int]bool{}
	for bar := 0; bar < 4; bar++ {
		chord := character.progression[bar]
		notes := bassNotes(float64(bar*4), bar, 0, 23, 48+chord.rootOffset, chord.third, 1, character)
		if len(notes) != 2 || notes[0].Beat == notes[1].Beat {
			t.Fatalf("energetic bass lacks a distinct response in bar %d: %+v", bar, notes)
		}
		intervals[notes[1].MIDI-notes[0].MIDI] = true
	}
	if len(intervals) < 3 {
		t.Fatalf("bass responses repeat one interval: %v", intervals)
	}
	calm := musicalStyles["calm"]
	if len(bassNotes(0, 0, 0, 23, 48, 4, 1, calm)) != 1 || len(bassNotes(36, 9, 2, 23, 53, 4, 1, calm)) != 2 {
		t.Fatal("calm bass did not develop between opening and bridge")
	}
}

// TestCounterlineUsesItsOwnRangeAndOnsets checks that the middle instrument answers between bass and lead while the opening leaves space.
func TestCounterlineUsesItsOwnRangeAndOnsets(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.7}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	totals := map[string]int{}
	counts := map[string]int{}
	for _, note := range plan.Notes {
		if note.Beat < 32 || note.Beat >= 48 {
			continue
		}
		if note.Voice == "bass" || note.Voice == "arp" || note.Voice == "lead" {
			totals[note.Voice] += note.MIDI
			counts[note.Voice]++
		}
	}
	if counts["bass"] == 0 || counts["arp"] == 0 || counts["lead"] == 0 || totals["bass"]/counts["bass"] >= totals["arp"]/counts["arp"] || totals["arp"]/counts["arp"] >= totals["lead"]/counts["lead"] {
		t.Fatalf("instrumental ranges collapsed together: totals=%v counts=%v", totals, counts)
	}
	if countVoiceInSection(plan, "pad", 0, 16) >= countVoiceInSection(plan, "pad", 32, 48) {
		t.Fatal("opening accompaniment does not leave space for the melody")
	}
}

// TestMelodyDevelopsWithoutRegisterJumps checks that the reprise changes the theme while consecutive pitches stay nearby.
func TestMelodyDevelopsWithoutRegisterJumps(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.8}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	opening := []int{}
	reprise := []int{}
	previous := 0
	for _, note := range plan.Notes {
		if note.Voice != "lead" {
			continue
		}
		if previous != 0 && (note.MIDI-previous > 12 || previous-note.MIDI > 12) {
			t.Fatalf("melody jumped too far from %d to %d", previous, note.MIDI)
		}
		previous = note.MIDI
		if note.Beat < 16 {
			opening = append(opening, note.MIDI)
		} else if note.Beat < 32 {
			reprise = append(reprise, note.MIDI)
		}
	}
	if len(opening) == 0 || len(opening) != len(reprise) || reflect.DeepEqual(opening, reprise) {
		t.Fatalf("the reprise did not develop the opening motif: %v / %v", opening, reprise)
	}
}

// TestMotifAndMelodicAnswerKeepTheirShape checks that a pentatonic peak does not wrap, quiet notes use the motif, and phrase endings land on their chord.
func TestMotifAndMelodicAnswerKeepTheirShape(t *testing.T) {
	var digest [32]byte
	digest[1], digest[2] = 1, 1
	motif := makeMotif(digest, 5)
	if motif[0] != 0 || motif[1] >= motif[2] || motif[3] != motif[1] {
		t.Fatalf("motif lost its arch: %v", motif)
	}
	character := musicalStyles["calm"]
	chord := chordSpec{rootOffset: 5, third: 4}
	notes, _ := melodyNotes(0, 0, 0, 23, 48, chord, motif, 0, 1, 0.1, character)
	if len(notes) != 2 || !pitchBelongsToChord(notes[1].MIDI, chordPitchClasses(53, 4)) {
		t.Fatalf("quiet melody does not settle into its chord: %+v", notes)
	}
	answer, _ := melodyNotes(12, 3, 0, 23, 48, chord, motif, notes[1].MIDI, 1, 0.1, character)
	if len(answer) != 1 || answer[0].MIDI%12 != (48+5+7)%12 {
		t.Fatalf("phrase ending missed the chord fifth: %+v", answer)
	}
}

// TestQuietMelodyChangesPitch checks that a sparse score still has distinct image-derived lead pitches rather than only chord roots.
func TestQuietMelodyChangesPitch(t *testing.T) {
	var firstDigest, secondDigest [32]byte
	secondDigest[1], secondDigest[2] = 1, 1
	first := composeMusic(analysis.StructuralAnalysis{Complexity: 0.1}, semantic.SemanticAnalysis{Description: "A serene landscape."}, firstDigest)
	second := composeMusic(analysis.StructuralAnalysis{Complexity: 0.1}, semantic.SemanticAnalysis{Description: "A serene landscape."}, secondDigest)
	firstLead := []int{}
	secondLead := []int{}
	seen := map[int]bool{}
	for _, note := range first.Notes {
		if note.Voice == "lead" && note.Beat < 16 {
			firstLead = append(firstLead, note.MIDI)
			seen[note.MIDI] = true
		}
	}
	for _, note := range second.Notes {
		if note.Voice == "lead" && note.Beat < 16 {
			secondLead = append(secondLead, note.MIDI)
		}
	}
	if len(seen) < 4 || reflect.DeepEqual(firstLead, secondLead) {
		t.Fatalf("sparse melody repeats too few pitches or ignores the motif: %v / %v", firstLead, secondLead)
	}
}

// TestLeadHasIndependentRegisterAndPhraseVariety checks that calm and energetic scenes produce clearly moving lead lines above the accompaniment.
func TestLeadHasIndependentRegisterAndPhraseVariety(t *testing.T) {
	for _, scene := range []string{"A serene landscape.", "A lively dancing scene."} {
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.2}, semantic.SemanticAnalysis{Description: scene}, [32]byte{})
		opening := map[int]bool{}
		previous := -1
		for _, note := range plan.Notes {
			if note.Voice != "lead" {
				continue
			}
			if note.MIDI < 67 || note.MIDI > 88 {
				t.Fatalf("lead left its own register for %q: %d", scene, note.MIDI)
			}
			if note.Beat < 16 {
				opening[note.MIDI] = true
			}
			if note.MIDI == previous && note.Beat < float64(plan.BeatsPerLoop-4) {
				t.Fatalf("lead repeated one pitch without a melodic step for %q", scene)
			}
			previous = note.MIDI
		}
		if len(opening) < 4 {
			t.Fatalf("opening phrase lacks pitch variety for %q: %v", scene, opening)
		}
	}
}

// TestOpeningMelodyAvoidsOneNoteDominance checks the quietest supported scene across tonal centers and motif variants.
func TestOpeningMelodyAvoidsOneNoteDominance(t *testing.T) {
	for center := byte(0); center < 12; center++ {
		for variant := byte(0); variant < 4; variant++ {
			var digest [32]byte
			digest[0], digest[1], digest[2] = center, variant&1, variant>>1
			plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.1}, semantic.SemanticAnalysis{Description: "A serene landscape."}, digest)
			counts := map[int]int{}
			total := 0
			for _, note := range plan.Notes {
				if note.Voice == "lead" && note.Beat < 16 {
					counts[note.MIDI]++
					total++
				}
			}
			for pitch, count := range counts {
				if count*2 >= total {
					t.Fatalf("one pitch dominates the opening: center=%d variant=%d pitch=%d count=%d/%d", center, variant, pitch, count, total)
				}
			}
		}
	}
}

// TestBridgeIntroducesSecondTheme checks that the B section has a distinct lead contour instead of repeating the opening unchanged.
func TestBridgeIntroducesSecondTheme(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.7}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	opening := []int{}
	bridge := []int{}
	for _, note := range plan.Notes {
		if note.Voice != "lead" {
			continue
		}
		if note.Beat < 16 {
			opening = append(opening, note.MIDI)
		} else if note.Beat >= 32 && note.Beat < 48 {
			bridge = append(bridge, note.MIDI)
		}
	}
	if len(opening) < 4 || len(bridge) < 4 || reflect.DeepEqual(opening[:4], bridge[:4]) {
		t.Fatalf("B theme repeats the first phrase: opening=%v bridge=%v", opening, bridge)
	}
}

// TestPhraseContoursTransformTheTheme checks that the second statement and bridge transform the motif while the return recalls it.
func TestPhraseContoursTransformTheTheme(t *testing.T) {
	motif := [4]int{0, 2, 4, 2}
	first := melodyDegree(motif, 0, 0, 1, 5)
	second := melodyDegree(motif, 4, 1, 1, 5)
	returning := melodyDegree(motif, 16, 4, 1, 5)
	if first == second || first != returning {
		t.Fatalf("motif lacks a recognizable variation: %d / %d / %d", first, second, returning)
	}
	forward := melodyDegree(motif, 8, 0, 1, 5)
	bridge := melodyDegree(motif, 8, 2, 1, 5)
	if forward == bridge {
		t.Fatalf("bridge did not reverse the motif: %d / %d", forward, bridge)
	}
}

// TestMelodyRhythmFollowsPhraseArc checks that a detailed image grows in the bridge and thins before the final cadence.
func TestMelodyRhythmFollowsPhraseArc(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.8}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	if countVoiceInSection(plan, "lead", 32, 48) <= countVoiceInSection(plan, "lead", 16, 32) {
		t.Fatal("the bridge should develop the melodic rhythm")
	}
	if countVoiceInSection(plan, "lead", float64(plan.BeatsPerLoop-16), float64(plan.BeatsPerLoop)) >= countVoiceInSection(plan, "lead", 32, 48) {
		t.Fatal("the conclusion should leave more space than the bridge")
	}
}

// TestPhraseRhythmDevelopsAndReturns checks that the second statement changes timing while the reprise recalls the opening rhythm.
func TestPhraseRhythmDevelopsAndReturns(t *testing.T) {
	character := musicalStyles["energetic"]
	opening := melodyRhythmForBar(0, 0, 0.6, character)
	variation := melodyRhythmForBar(4, 1, 0.6, character)
	returning := melodyRhythmForBar(16, 4, 0.6, character)
	if reflect.DeepEqual(opening, variation) || !reflect.DeepEqual(opening, returning) {
		t.Fatalf("phrase rhythm lacks a varied return: %v / %v / %v", opening, variation, returning)
	}
	for bar := 0; bar < len(formEpisodes)*4; bar++ {
		rhythm := melodyRhythmForBar(bar, formEpisodes[bar/4].role, 0.6, character)
		for index, beat := range rhythm {
			if beat < 0 || beat >= 4 || (index > 0 && beat <= rhythm[index-1]) {
				t.Fatalf("invalid rhythm at bar %d: %v", bar, rhythm)
			}
		}
	}
}

// TestPhraseDynamicsBreathe checks that the accompaniment grows inside a phrase and releases before its cadence.
func TestPhraseDynamicsBreathe(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.6}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	pads := map[int]MusicNote{}
	for _, note := range plan.Notes {
		if note.Voice == "pad" && note.Beat == float64(int(note.Beat)) {
			bar := int(note.Beat) / 4
			if _, exists := pads[bar]; !exists {
				pads[bar] = note
			}
		}
	}
	if pads[2].Velocity <= pads[0].Velocity || pads[3].Velocity >= pads[2].Velocity {
		t.Fatalf("opening phrase has no intensity arc: %+v %+v %+v", pads[0], pads[2], pads[3])
	}
	if pads[1].Duration >= pads[0].Duration || pads[3].Duration >= pads[2].Duration {
		t.Fatalf("accompaniment does not leave phrase space: %+v %+v %+v %+v", pads[0], pads[1], pads[2], pads[3])
	}
	lastBar := plan.BeatsPerLoop/4 - 1
	if pads[lastBar].Velocity >= pads[lastBar-3].Velocity {
		t.Fatalf("ending does not fade: %+v %+v", pads[lastBar-3], pads[lastBar])
	}
}

// TestMusicFormBreaksAndReturns checks that the interlude makes space before the returning theme regains its pulse.
func TestMusicFormBreaksAndReturns(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.7}, semantic.SemanticAnalysis{Description: "A lively dancing scene."}, [32]byte{})
	if plan.Sections[3].Label != "Interlude" || plan.Sections[4].Label != "Retour" {
		t.Fatalf("missing form contrast: %+v", plan.Sections)
	}
	if countVoiceInSection(plan, "pulse", 48, 64) != 0 || countVoiceInSection(plan, "lead", 48, 64) >= countVoiceInSection(plan, "lead", 64, 80) || countVoiceInSection(plan, "pulse", 64, 80) == 0 {
		t.Fatal("the interlude and return do not contrast")
	}
	if len(plan.Notes) > 1200 {
		t.Fatalf("score exceeds the player event limit: %d", len(plan.Notes))
	}
}

// TestHummingDependsOnScene checks that the vocal response is sparse and appears only for suitable interpreted artwork.
func TestHummingDependsOnScene(t *testing.T) {
	quiet := analysis.StructuralAnalysis{Complexity: 0.2}
	calm := composeMusic(quiet, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	person := composeMusic(quiet, semantic.SemanticAnalysis{Description: "A drawing of a stick figure."}, [32]byte{})
	threat := composeMusic(quiet, semantic.SemanticAnalysis{Description: "A threatening violent scene."}, [32]byte{})
	busy := composeMusic(analysis.StructuralAnalysis{Complexity: 0.8}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	if countVoiceInSection(calm, "hum", 48, 64) != 2 || countVoiceInSection(person, "hum", 48, 64) != 2 {
		t.Fatal("suitable calm or human scenes need a sparse vocal response")
	}
	if countVoiceInSection(threat, "hum", 0, float64(threat.BeatsPerLoop)) != 0 || countVoiceInSection(busy, "hum", 0, float64(busy.BeatsPerLoop)) != 0 {
		t.Fatal("vocal response should be omitted for tense or complex scenes")
	}
}

// TestExpressNotesPreservesAnchors checks that humanized melody timing stays in its bar while sustained chords remain aligned.
func TestExpressNotesPreservesAnchors(t *testing.T) {
	notes := []MusicNote{
		{Beat: 0, Velocity: 0.3, Voice: "pad"},
		{Beat: 0, Velocity: 0.6, Voice: "lead"},
		{Beat: 3.99, Velocity: 0.4, Voice: "arp"},
		{Beat: 92, Velocity: 0.6, Voice: "lead"},
	}
	var digest [32]byte
	for index := range digest {
		digest[index] = 255
	}
	expressNotes(notes, digest, 92)
	if notes[0].Beat != 0 || notes[0].Velocity != 0.3 || notes[1].Beat <= 0 || notes[2].Beat >= 4 || notes[3].Beat != 92 {
		t.Fatalf("expression moved an anchor or crossed a bar: %+v", notes)
	}
}

// TestMusicServiceComposeRejectsInvalidImage checks that the score cannot be generated from an unbounded or malformed image.
func TestMusicServiceComposeRejectsInvalidImage(t *testing.T) {
	if _, err := (&MusicService{}).Compose("data:image/png;base64,AA==", semantic.SemanticAnalysis{}); err == nil {
		t.Fatal("expected malformed image to be rejected")
	}
}
