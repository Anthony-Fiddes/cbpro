package cbpro

const approxNumProducts = 250

type Product struct {
	BaseCurrency    string `json:"base_currency"`
	BaseIncrement   string `json:"base_increment"`
	BaseMaxSize     string `json:"base_max_size"`
	BaseMinSize     string `json:"base_min_size"`
	CancelOnly      bool   `json:"cancel_only"`
	DisplayName     string `json:"display_name"`
	ID              string `json:"id"`
	LimitOnly       bool   `json:"limit_only"`
	MaxMarketFunds  string `json:"max_market_funds"`
	MinMarketFunds  string `json:"min_market_funds"`
	PostOnly        bool   `json:"post_only"`
	QuoteCurrency   string `json:"quote_currency"`
	QuoteIncrement  string `json:"quote_increment"`
	Status          string `json:"status"`
	StatusMessage   string `json:"status_message"`
	TradingDisabled bool   `json:"trading_disabled"`
}
