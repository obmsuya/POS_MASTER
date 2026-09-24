package i18n

import "fmt"

type Lang string

const (
	Swahili Lang = "sw"
	English Lang = "en"
)

var current Lang = Swahili

func Set(lang Lang) {
	current = lang
}

var strings = map[string]map[Lang]string{
	"pick_language":       {Swahili: "Chagua lugha", English: "Choose a language"},
	"welcome":             {Swahili: "Karibu, %s!", English: "Welcome, %s!"},
	"loading":             {Swahili: "Subiri kidogo...", English: "One moment..."},
	"connecting":          {Swahili: "Inaunganisha na seva...", English: "Connecting to the server..."},
	"phone_label":         {Swahili: "Namba ya simu", English: "Phone number"},
	"password_label":      {Swahili: "Nywila", English: "Password"},
	"login_title":         {Swahili: "Ingia FALTASI POS", English: "Sign in to FALTASI POS"},
	"login_failed":        {Swahili: "Imeshindwa kuingia: %s", English: "Sign in failed: %s"},
	"invalid_phone":       {Swahili: "Namba ya simu si sahihi. Tumia mfano 07XXXXXXXX.", English: "That phone number isn't valid. Use a format like 07XXXXXXXX."},
	"invalid_hwid":        {Swahili: "Umekosea kidogo — hii ID ya POS (Hardware ID) si sahihi.", English: "That doesn't look right — check the Hardware ID."},
	"invalid_amount":      {Swahili: "Kiasi kilicholipwa si sahihi. Lazima kiwe zaidi ya sifuri.", English: "That amount isn't valid. It must be greater than zero."},
	"invalid_package":     {Swahili: "Chagua kifurushi kwenye orodha.", English: "Pick a package from the list."},
	"role_not_allowed":    {Swahili: "Akaunti hii haiwezi kutumia FALTASI POS.", English: "This account can't use FALTASI POS."},
	"main_menu_title":     {Swahili: "Dashibodi Kuu", English: "Main Menu"},
	"menu_activate":       {Swahili: "Anzisha au Ongeza Leseni", English: "Activate or Renew License"},
	"menu_lookup":         {Swahili: "Tafuta Mteja", English: "Search Customer"},
	"menu_history":        {Swahili: "Historia ya Malipo", English: "Payment History"},
	"menu_packages":       {Swahili: "Simamia Vifurushi", English: "Manage Packages"},
	"menu_exit":           {Swahili: "Toka", English: "Exit"},
	"lookup_prompt":       {Swahili: "Weka ID ya POS (Hardware ID) ya mteja", English: "Enter the customer's Hardware ID"},
	"lookup_or_phone":     {Swahili: "Au weka namba ya simu (kama hujui ID)", English: "Or enter their phone number (if you don't have the ID)"},
	"customer_found":      {Swahili: "Mteja amepatikana", English: "Customer found"},
	"customer_not_found":  {Swahili: "Mteja huyu hajapatikana — tutamsajili kama mpya.", English: "No existing customer found — we'll register them as new."},
	"new_customer_phone":  {Swahili: "Weka namba ya simu ya mteja mpya", English: "Enter the new customer's phone number"},
	"pick_package":        {Swahili: "Chagua kifurushi", English: "Choose a package"},
	"amount_label":        {Swahili: "Kiasi kilicholipwa (TZS)", English: "Amount paid (TZS)"},
	"method_label":        {Swahili: "Njia ya malipo", English: "Payment method"},
	"method_cash":         {Swahili: "Fedha taslimu", English: "Cash"},
	"method_mobile_money": {Swahili: "Pesa za Simu", English: "Mobile Money"},
	"method_bank":         {Swahili: "Uhamisho wa Benki", English: "Bank Transfer"},
	"reference_label":     {Swahili: "Kumbukumbu (hiari) — mf. namba ya muamala", English: "Reference (optional) — e.g. transaction code"},
	"confirm_title":       {Swahili: "Thibitisha Malipo", English: "Confirm Payment"},
	"confirm_prompt":      {Swahili: "Endelea na malipo haya?", English: "Proceed with this payment?"},
	"activating":          {Swahili: "Inawasha leseni...", English: "Activating license..."},
	"activated_title":     {Swahili: "Imefanikiwa!", English: "Success!"},
	"activated_body":      {Swahili: "%s — siku %d zimeongezwa", English: "%s — %d days added"},
	"expires_label":       {Swahili: "Inaisha tarehe", English: "Expires"},
	"recorded_by":         {Swahili: "Imesajiliwa na", English: "Recorded by"},
	"press_enter":         {Swahili: "Bonyeza Enter kuendelea", English: "Press Enter to continue"},
	"error_generic":       {Swahili: "Hitilafu: %s", English: "Error: %s"},
	"error_network":       {Swahili: "Imeshindwa kufikia seva. Hakikisha una mtandao.", English: "Couldn't reach the server. Check your connection."},
	"history_title":       {Swahili: "Malipo ya Hivi Karibuni", English: "Recent Payments"},
	"history_empty":       {Swahili: "Hakuna malipo bado.", English: "No payments recorded yet."},
	"packages_title":      {Swahili: "Simamia Vifurushi", English: "Manage Packages"},
	"pkg_create":          {Swahili: "Unda Kifurushi Kipya", English: "Create New Package"},
	"pkg_edit":            {Swahili: "Hariri Kifurushi", English: "Edit a Package"},
	"pkg_back":            {Swahili: "Rudi", English: "Back"},
	"pkg_name_label":      {Swahili: "Jina la kifurushi", English: "Package name"},
	"pkg_price_label":     {Swahili: "Bei (TZS)", English: "Price (TZS)"},
	"pkg_days_label":      {Swahili: "Siku ngapi", English: "Days granted"},
	"pkg_devices_label":   {Swahili: "Vifaa vingapi vinaruhusiwa", English: "Max devices allowed"},
	"pkg_edit_notice":     {Swahili: "Kumbuka: mabadiliko haya hayamuathiri mteja aliyeshanunua kwa bei ya zamani — ni kwa ajili ya wateja wapya tu.", English: "Note: this only affects new purchases — customers already on this package keep what they paid for."},
	"pkg_created":         {Swahili: "Kifurushi kimeundwa!", English: "Package created!"},
	"pkg_updated":         {Swahili: "Kifurushi kimesasishwa!", English: "Package updated!"},
	"goodbye":             {Swahili: "Asante kwa kutumia FALTASI POS. Kwaheri!", English: "Thank you for using FALTASI POS. Goodbye!"},
	"invalid_input_retry": {Swahili: "Umekosea kidogo, jaribu tena.", English: "That's not quite right — try again."},
	"required_field":      {Swahili: "Sehemu hii inahitajika.", English: "This field is required."},
	"session_restored":    {Swahili: "Karibu tena, %s.", English: "Welcome back, %s."},
	"update_available":    {Swahili: "Toleo jipya %s linapatikana. Sasisha sasa?", English: "Version %s is available. Update now?"},
	"updating":            {Swahili: "Inasasisha...", English: "Updating..."},
	"update_done":         {Swahili: "Imesasishwa! Fungua tena programu kutumia toleo jipya.", English: "Updated! Restart the app to use the new version."},
	"update_failed":       {Swahili: "Imeshindwa kusasisha: %s", English: "Update failed: %s"},
	"days_suffix":         {Swahili: "siku", English: "days"},
	"devices_suffix":      {Swahili: "vifaa", English: "devices"},
	"no_packages":         {Swahili: "Hakuna vifurushi vilivyopo. Muulize msimamizi aunde kimoja.", English: "No packages exist yet. Ask an admin to create one."},
}

func T(key string, args ...interface{}) string {
	entry, ok := strings[key]
	if !ok {
		return key
	}
	template, ok := entry[current]
	if !ok {
		template = entry[English]
	}
	if len(args) == 0 {
		return template
	}
	return fmt.Sprintf(template, args...)
}
