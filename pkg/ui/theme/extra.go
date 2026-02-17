package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func animatedTick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func cardInnerWidth(cardWidth int) int {
	innerWidth := cardWidth - 6
	if innerWidth < 1 {
		innerWidth = 1
	}
	if innerWidth > cardWidth {
		innerWidth = cardWidth
	}
	return innerWidth
}

func appendWrappedLines(lines []string, prefix, text string, width, maxLines int) []string {
	if width <= 0 || maxLines <= 0 {
		return lines
	}
	return append(lines, wrapAndClamp(prefix, text, width, maxLines)...)
}

func appendRawLines(lines []string, block string, width, maxLines int) []string {
	if width <= 0 || maxLines <= 0 {
		return lines
	}
	normalized := strings.ReplaceAll(strings.ReplaceAll(block, "\r\n", "\n"), "\r", "\n")
	raw := strings.Split(normalized, "\n")
	return append(lines, clampLines(raw, maxLines, width)...)
}

func buildCardBody(lines []string, width, maxLines int) string {
	if width <= 0 {
		return ""
	}
	if maxLines <= 0 {
		maxLines = 1
	}
	return strings.Join(clampLines(lines, maxLines, width), "\n")
}

type AuroraGridTheme struct {
	frame int
}

func NewAuroraGridTheme() Theme { return &AuroraGridTheme{} }
func (t *AuroraGridTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 140)
}
func (t *AuroraGridTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 140)
	}
	return t, nil
}
func (t *AuroraGridTheme) Name() string        { return "Aurora Grid" }
func (t *AuroraGridTheme) Description() string { return "Neon grid with shifting aurora edge glow" }
func (t *AuroraGridTheme) IsCompatible(c caps.Capabilities) bool {
	// Requires at least basic color support; otherwise the theme loses most of its effect.
	return c.ColorProfile >= termenv.ANSI
}
func (t *AuroraGridTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	colors := []string{"#00FFD5", "#00B8FF", "#7AFF00", "#00FF7A"}
	pulse := lipgloss.Color(colors[t.frame%len(colors)])
	boxW := clamp(width-10, 36, 82)
	gridW := clamp(boxW-4, 20, 78)
	var bg strings.Builder
	for y := 0; y < 5; y++ {
		line := strings.Repeat(". ", gridW/2)
		if (y+t.frame)%2 == 0 {
			line = strings.Repeat(" : ", gridW/3)
		}
		bg.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#114444")).Render(line))
		if y < 4 {
			bg.WriteRune('\n')
		}
	}

	innerW := cardInnerWidth(boxW)
	maxBodyLines := remainingRows(height, 9, 8)

	lines := []string{fmt.Sprintf("AURORA NODE %02d", q.ID), ""}
	lines = appendWrappedLines(lines, "", q.Text, innerW, 4)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "INPUT > ", inputView, innerW, 2)
	if hint != "" {
		lines = append(lines, "")
		lines = appendWrappedLines(lines, "HINT: ", hint, innerW, 2)
	}
	body := buildCardBody(lines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(pulse).
		Foreground(lipgloss.Color("#D6FFF6")).
		Background(lipgloss.Color("#021B1D")).
		Padding(1).
		Width(boxW).
		Render(body)

	stack := lipgloss.JoinVertical(lipgloss.Left, bg.String(), "", card)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, stack)
}

type RadarSweepTheme struct {
	frame int
}

func NewRadarSweepTheme() Theme { return &RadarSweepTheme{} }
func (t *RadarSweepTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 120)
}
func (t *RadarSweepTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 120)
	}
	return t, nil
}
func (t *RadarSweepTheme) Name() string        { return "Radar Sweep" }
func (t *RadarSweepTheme) Description() string { return "Military radar panel with moving sweep bar" }
func (t *RadarSweepTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile >= termenv.ANSI
}
func (t *RadarSweepTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	panelW := clamp(width-12, 34, 76)
	innerW := cardInnerWidth(panelW)
	sweepW := clamp(innerW-4, 8, 64)
	pos := 0
	if sweepW > 0 {
		pos = t.frame % sweepW
	}
	barPlain := strings.Repeat(".", pos) + "|" + strings.Repeat(".", sweepW-pos-1)

	status := []string{"LOCKING", "TRACKING", "STABILIZING"}
	head := fmt.Sprintf("RADAR %s", status[t.frame%len(status)])

	maxBodyLines := remainingRows(height, 6, 8)
	lines := []string{head, "[" + barPlain + "]", "", fmt.Sprintf("TARGET PROMPT %02d", q.ID)}
	lines = appendWrappedLines(lines, "", q.Text, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "CMD > ", inputView, innerW, 2)
	if hint != "" {
		lines = append(lines, "")
		lines = appendWrappedLines(lines, "AUX: ", hint, innerW, 2)
	}
	body := buildCardBody(lines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00AA44")).
		Background(lipgloss.Color("#031106")).
		Foreground(lipgloss.Color("#A4FFBF")).
		Padding(1).
		Width(panelW).
		Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type BlueprintTheme struct {
	frame int
}

func NewBlueprintTheme() Theme { return &BlueprintTheme{} }
func (t *BlueprintTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 180)
}
func (t *BlueprintTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 180)
	}
	return t, nil
}
func (t *BlueprintTheme) Name() string        { return "Blueprint Ops" }
func (t *BlueprintTheme) Description() string { return "Engineering blueprint with active node blink" }
func (t *BlueprintTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile >= termenv.ANSI
}
func (t *BlueprintTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	nodes := []string{"[A]", "[B]", "[C]", "[D]"}
	active := t.frame % len(nodes)
	nodes[active] = "<" + string(rune('A'+active)) + ">"
	diagram := fmt.Sprintf("%s---%s\n |   |\n%s---%s", nodes[0], nodes[1], nodes[2], nodes[3])

	cardW := clamp(width-10, 36, 78)
	innerW := cardInnerWidth(cardW)
	maxBodyLines := remainingRows(height, 6, 8)

	lines := []string{fmt.Sprintf("BLUEPRINT LAYER %02d", q.ID), ""}
	lines = appendRawLines(lines, diagram, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "", q.Text, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "INPUT: ", inputView, innerW, 2)
	if hint != "" {
		lines = append(lines, "")
		lines = appendWrappedLines(lines, "NOTE: ", hint, innerW, 2)
	}
	body := buildCardBody(lines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#65B2FF")).
		Background(lipgloss.Color("#06172B")).
		Foreground(lipgloss.Color("#BFE2FF")).
		Padding(1, 2).
		Width(cardW).
		Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type GlitchLabTheme struct {
	frame int
}

func NewGlitchLabTheme() Theme { return &GlitchLabTheme{} }
func (t *GlitchLabTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 100)
}
func (t *GlitchLabTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 100)
	}
	return t, nil
}
func (t *GlitchLabTheme) Name() string { return "Glitch Lab" }
func (t *GlitchLabTheme) Description() string {
	return "Corrupted terminal with controlled text distortion"
}
func (t *GlitchLabTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile >= termenv.ANSI
}
func (t *GlitchLabTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	r := []rune(q.Text)
	symbols := []rune{'#', '%', '&', '?'}
	for i := 0; i < len(r); i++ {
		if r[i] == ' ' {
			continue
		}
		if (i+t.frame)%11 == 0 {
			r[i] = symbols[(i+t.frame)%len(symbols)]
		}
	}
	glitched := string(r)
	banner := "GLITCH LAB :: SIGNAL INTEGRITY " + fmt.Sprintf("%d%%", 80+(t.frame%20))

	cardW := clamp(width-10, 38, 82)
	innerW := cardInnerWidth(cardW)
	maxBodyLines := remainingRows(height, 6, 8)

	lines := []string{banner, ""}
	lines = appendWrappedLines(lines, "", glitched, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "RAW: ", q.Text, innerW, 2)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "INJECT > ", inputView, innerW, 2)
	if hint != "" {
		lines = append(lines, "")
		lines = appendWrappedLines(lines, "RECOVERY HINT: ", hint, innerW, 2)
	}
	body := buildCardBody(lines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#FF2255")).
		Background(lipgloss.Color("#12060A")).
		Foreground(lipgloss.Color("#FFD0D9")).
		Padding(1).
		Width(cardW).
		Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type VaultLedgerTheme struct {
	frame int
}

func NewVaultLedgerTheme() Theme { return &VaultLedgerTheme{} }
func (t *VaultLedgerTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 200)
}
func (t *VaultLedgerTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 200)
	}
	return t, nil
}
func (t *VaultLedgerTheme) Name() string { return "Vault Ledger" }
func (t *VaultLedgerTheme) Description() string {
	return "Brass vault interface with rotating dial indicator"
}
func (t *VaultLedgerTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile >= termenv.ANSI
}
func (t *VaultLedgerTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	dial := []string{"|", "/", "-", "\\"}
	code := fmt.Sprintf("%02d-%02d-%02d", (q.ID+t.frame)%100, (q.ID*3+t.frame)%100, (q.ID*7+t.frame)%100)

	cardW := clamp(width-10, 38, 82)
	innerW := cardInnerWidth(cardW)
	maxBodyLines := remainingRows(height, 6, 8)

	lines := []string{
		fmt.Sprintf("VAULT LEDGER %s", time.Now().Format("15:04:05")),
		"",
		fmt.Sprintf("DIAL: %s", dial[t.frame%len(dial)]),
		fmt.Sprintf("CODE TRACK: %s", code),
		"",
	}
	lines = appendWrappedLines(lines, "", q.Text, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "KEY > ", inputView, innerW, 2)
	if hint != "" {
		lines = append(lines, "")
		lines = appendWrappedLines(lines, "CLUE: ", hint, innerW, 2)
	}
	body := buildCardBody(lines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#D8B065")).
		Background(lipgloss.Color("#1E1608")).
		Foreground(lipgloss.Color("#FFE4A8")).
		Padding(1).
		Width(cardW).
		Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type SonarTheme struct {
	frame int
}

func NewSonarTheme() Theme { return &SonarTheme{} }
func (t *SonarTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 130)
}
func (t *SonarTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 130)
	}
	return t, nil
}
func (t *SonarTheme) Name() string        { return "Deep Sonar" }
func (t *SonarTheme) Description() string { return "Subsea sonar display with rolling pulse lines" }
func (t *SonarTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile >= termenv.ANSI
}
func (t *SonarTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	cardW := clamp(width-10, 36, 80)
	innerW := cardInnerWidth(cardW)
	w := innerW
	if w < 12 {
		w = 12
	}
	p := t.frame % w
	lines := []string{
		strings.Repeat("~", w),
		strings.Repeat(" ", p) + "~~~" + strings.Repeat(" ", clamp(w-p-3, 0, w)),
		strings.Repeat(" ", (p*2)%w) + "~~~~",
	}
	waves := strings.Join(lines, "\n")

	maxBodyLines := remainingRows(height, 6, 8)
	bodyLines := []string{"DEEP SONAR", ""}
	bodyLines = appendRawLines(bodyLines, waves, innerW, 3)
	bodyLines = append(bodyLines, "", fmt.Sprintf("PING %02d", t.frame%60))
	bodyLines = appendWrappedLines(bodyLines, "", q.Text, innerW, 2)
	bodyLines = append(bodyLines, "")
	bodyLines = appendWrappedLines(bodyLines, "CMD ", inputView, innerW, 2)
	if hint != "" {
		bodyLines = append(bodyLines, "")
		bodyLines = appendWrappedLines(bodyLines, "ECHO: ", hint, innerW, 2)
	}
	body := buildCardBody(bodyLines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#33C4FF")).
		Background(lipgloss.Color("#031624")).
		Foreground(lipgloss.Color("#B8EBFF")).
		Padding(1).
		Width(cardW).
		Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type EmberForgeTheme struct {
	frame int
}

func NewEmberForgeTheme() Theme { return &EmberForgeTheme{} }
func (t *EmberForgeTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 110)
}
func (t *EmberForgeTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 110)
	}
	return t, nil
}
func (t *EmberForgeTheme) Name() string        { return "Ember Forge" }
func (t *EmberForgeTheme) Description() string { return "Industrial forge terminal with rising sparks" }
func (t *EmberForgeTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile >= termenv.ANSI
}
func (t *EmberForgeTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	sparks := []string{" .  *   . ", "   *  .   ", " *   .   *", "  . *   . "}
	sparkBand := sparks[t.frame%len(sparks)] + sparks[(t.frame+1)%len(sparks)] + sparks[(t.frame+2)%len(sparks)]

	cardW := clamp(width-10, 34, 76)
	innerW := cardInnerWidth(cardW)
	maxBodyLines := remainingRows(height, 6, 8)

	lines := []string{"EMBER FORGE", truncateToWidth(sparkBand, innerW), ""}
	lines = appendWrappedLines(lines, "", q.Text, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "HAMMER > ", inputView, innerW, 2)
	if hint != "" {
		lines = append(lines, "")
		lines = appendWrappedLines(lines, "COOLANT NOTE: ", hint, innerW, 2)
	}
	body := buildCardBody(lines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#FF7A2F")).
		Background(lipgloss.Color("#1A0D06")).
		Foreground(lipgloss.Color("#FFD1B0")).
		Padding(1).
		Width(cardW).
		Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type FrostbyteTheme struct {
	frame int
}

func NewFrostbyteTheme() Theme { return &FrostbyteTheme{} }
func (t *FrostbyteTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 170)
}
func (t *FrostbyteTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 170)
	}
	return t, nil
}
func (t *FrostbyteTheme) Name() string        { return "Frostbyte Core" }
func (t *FrostbyteTheme) Description() string { return "Cold system shell with pulsing freeze gauge" }
func (t *FrostbyteTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile >= termenv.ANSI
}
func (t *FrostbyteTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	gauge := strings.Repeat("#", (t.frame%10)+1) + strings.Repeat("-", 10-((t.frame%10)+1))

	cardW := clamp(width-10, 36, 80)
	innerW := cardInnerWidth(cardW)
	maxBodyLines := remainingRows(height, 6, 8)

	lines := []string{
		"FROSTBYTE CORE",
		fmt.Sprintf("THERMAL LOCK [%s]", gauge),
		"",
	}
	lines = appendWrappedLines(lines, "", q.Text, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "ICECMD > ", inputView, innerW, 2)
	if hint != "" {
		lines = append(lines, "")
		lines = appendWrappedLines(lines, "DEICE: ", hint, innerW, 2)
	}
	body := buildCardBody(lines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#9FD8FF")).
		Background(lipgloss.Color("#07131C")).
		Foreground(lipgloss.Color("#E3F6FF")).
		Padding(1).
		Width(cardW).
		Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type NoirDossierTheme struct {
	frame int
}

func NewNoirDossierTheme() Theme { return &NoirDossierTheme{} }
func (t *NoirDossierTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 160)
}
func (t *NoirDossierTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 160)
	}
	return t, nil
}
func (t *NoirDossierTheme) Name() string { return "Noir Dossier" }
func (t *NoirDossierTheme) Description() string {
	return "Monochrome dossier view with rolling scanline marker"
}
func (t *NoirDossierTheme) IsCompatible(c caps.Capabilities) bool {
	return true
}
func (t *NoirDossierTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	marker := []string{"[ ]", "[=]", "[#]", "[=]"}[t.frame%4]

	cardW := clamp(width-10, 36, 78)
	innerW := cardInnerWidth(cardW)
	maxBodyLines := remainingRows(height, 6, 8)

	lines := []string{fmt.Sprintf("DOSSIER FILE %02d %s", q.ID, marker), "", "SUBJECT PROMPT:"}
	lines = appendWrappedLines(lines, "", q.Text, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "RESPONSE: ", inputView, innerW, 2)
	if hint != "" {
		lines = append(lines, "")
		lines = appendWrappedLines(lines, "NOTE: ", hint, innerW, 2)
	}
	body := buildCardBody(lines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#A0A0A0")).
		Foreground(lipgloss.Color("#E6E6E6")).
		Background(lipgloss.Color("#0E0E0E")).
		Padding(1, 2).
		Width(cardW).
		Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type CircuitBoardTheme struct {
	frame int
}

func NewCircuitBoardTheme() Theme { return &CircuitBoardTheme{} }
func (t *CircuitBoardTheme) Init() tea.Cmd {
	return animatedTick(time.Millisecond * 125)
}
func (t *CircuitBoardTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, animatedTick(time.Millisecond * 125)
	}
	return t, nil
}
func (t *CircuitBoardTheme) Name() string        { return "Circuit Board" }
func (t *CircuitBoardTheme) Description() string { return "PCB trace map with moving signal pulse" }
func (t *CircuitBoardTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile >= termenv.ANSI
}
func (t *CircuitBoardTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	path := []string{
		"o---+----+------o",
		"    |    |       ",
		"o---+----+---o   ",
	}
	signal := t.frame % len(path[0])
	for i := range path {
		r := []rune(path[i])
		if signal < len(r) && r[signal] != ' ' {
			r[signal] = '*'
		}
		path[i] = string(r)
	}
	trace := strings.Join(path, "\n")

	cardW := clamp(width-10, 36, 78)
	innerW := cardInnerWidth(cardW)
	maxBodyLines := remainingRows(height, 6, 8)

	lines := []string{fmt.Sprintf("PCB NET %02d", q.ID), ""}
	lines = appendRawLines(lines, trace, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "", q.Text, innerW, 3)
	lines = append(lines, "")
	lines = appendWrappedLines(lines, "BUS > ", inputView, innerW, 2)
	if hint != "" {
		lines = append(lines, "")
		lines = appendWrappedLines(lines, "TRACE NOTE: ", hint, innerW, 2)
	}
	body := buildCardBody(lines, innerW, maxBodyLines)

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5EEA8F")).
		Background(lipgloss.Color("#08180D")).
		Foreground(lipgloss.Color("#C8FFD9")).
		Padding(1).
		Width(cardW).
		Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

func init() {
	Register(NewAuroraGridTheme)
	Register(NewRadarSweepTheme)
	Register(NewBlueprintTheme)
	Register(NewGlitchLabTheme)
	Register(NewVaultLedgerTheme)
	Register(NewSonarTheme)
	Register(NewEmberForgeTheme)
	Register(NewFrostbyteTheme)
	Register(NewNoirDossierTheme)
	Register(NewCircuitBoardTheme)
}
