package main

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters"
	ackscanners "git.netsplit.it/enrico204/blocklists/tools/internal/filters/acknowledged_scanners"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/apacheconf"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/csv"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/dataplane"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/dshield"
	et_pix "git.netsplit.it/enrico204/blocklists/tools/internal/filters/emergingthreats-pix"
	gpfcomics "git.netsplit.it/enrico204/blocklists/tools/internal/filters/gpf_comics"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/hostdeny"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/ip9datacenters"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/local_files"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/p2p"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/plaintext"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters/torproject_exits"
	webpagerx "git.netsplit.it/enrico204/blocklists/tools/internal/filters/webpage_regex"
	"go.uber.org/zap"
	"regexp"
)

// registerFilters contains all the filter tags/names associated with an instance of filters.Filter
func registerFilters(logger *zap.SugaredLogger, cacheDir string) map[string]filters.Filter {
	logger.Info("initializing filters")
	var filters = make(map[string]filters.Filter)
	filters["acknowledged_scanners"] = ackscanners.New(cacheDir)
	filters["csv"] = csv.New(cacheDir, csv.Options{SkipHeader: true})
	filters["csv_semicolon"] = csv.New(cacheDir, csv.Options{SkipHeader: true, Comma: ';'})
	filters["p2p_gz"] = p2p.New(cacheDir)
	filters["dshield_parser"] = dshield.New(cacheDir)
	filters["hostdeny"] = hostdeny.New(cacheDir)
	filters["dataplane_column3"] = dataplane.New(cacheDir)
	filters["pix_deny_rules_to_ipv4"] = et_pix.New(cacheDir)
	filters["file"] = localfiles.New()
	filters["plaintext"] = plaintext.New(cacheDir, plaintext.Options{})
	filters["unzip_and_split_csv"] = plaintext.New(cacheDir, plaintext.Options{
		Unzip:              true,
		UseCustomSeparator: true,
		CustomSeparator:    ',',
	})
	filters["unzip_and_extract"] = plaintext.New(cacheDir, plaintext.Options{Unzip: true})
	filters["botscout_filter"] = webpagerx.New(cacheDir, regexp.MustCompile(`href="/ipcheck\.htm\?ip=([0-9a-fA-F:.]+)"`))
	filters["parse_graphiclineweb"] = webpagerx.New(cacheDir, regexp.MustCompile(`<td style="width:120px;">([0-9a-fA-F:.]+)</td>`))
	filters["parse_cleantalk"] = webpagerx.New(cacheDir, regexp.MustCompile(`<a href="https://cleantalk\.org/blacklists/([0-9a-fA-F:.]+)"`))
	filters["gpf_comics"] = gpfcomics.New(cacheDir)
	filters["snort_alert_rules_to_ipv4"] = webpagerx.New(cacheDir, regexp.MustCompile(`([0-9a-fA-F]+[:.][0-9a-fA-F:.]+[:.][0-9a-fA-F:.]+),`))
	filters["torproject_exits"] = torprojectexits.New(cacheDir)
	filters["parse_rss_proxy"] = webpagerx.New(cacheDir, regexp.MustCompile(`<prx:ip>([0-9a-fA-F:.]+)</prx:ip>`))
	filters["extract_ipv4_from_any_file"] = webpagerx.New(cacheDir, regexp.MustCompile(`(((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9]))`))
	filters["parse_maxmind_proxy_fraud"] = webpagerx.New(cacheDir, regexp.MustCompile(`href="/en/high-risk-ip-sample/([0-9a-fA-F:.]+)"`))

	filters["csv_comma_first_column"] = csv.New(cacheDir, csv.Options{SkipHeader: true, Column: 1, Comment: '#'})
	filters["parse_corpus"] = webpagerx.New(cacheDir, regexp.MustCompile(`<a href='https://www\.virustotal\.com/en/ip-address/([0-9a-fA-F:.]+)/information/'>`))
	filters["parse_cvs_dyndns_ponmocup"] = csv.New(cacheDir, csv.Options{SkipHeader: true, Column: 1})
	filters["parse_cta_cryptowall"] = csv.New(cacheDir, csv.Options{SkipHeader: true, Column: 2})

	filters["parse_client9_ipcat_datacenters"] = ip9datacenters.New(cacheDir)
	filters["parse_turris_greylist"] = csv.New(cacheDir, csv.Options{SkipHeader: true, Comment: '#'})
	filters["apacheconf"] = apacheconf.New(cacheDir)
	filters["parse_php_rss"] = webpagerx.New(cacheDir, regexp.MustCompile(`<title>([0-9a-fA-F:.]+) \|.*</title>`))
	filters["gz_second_word"] = csv.New(cacheDir, csv.Options{Column: 1, Comma: ' '})
	filters["urlvir_last"] = webpagerx.New(cacheDir, regexp.MustCompile(`<td>([0-9a-fA-F]+[:.][0-9a-fA-F:.]+)</td>`))

	// TODO: regex limited to IPv4
	filters["urlhaus"] = webpagerx.New(cacheDir, regexp.MustCompile(`https?://([0-9a-fA-F.]+)(:[0-9]+)?/`))

	filters["local_files"] = localfiles.New()

	filters["parse_pushing_inertia"] = filters["apacheconf"]
	filters["remove_comments_semi_colon"] = filters["plaintext"]
	filters["cat"] = filters["plaintext"]
	filters["remove_comments"] = filters["plaintext"]
	filters["p2p_gz_proxy"] = filters["p2p_gz"] // Here the original p2p_gz_proxy filters out all Tor: lines
	filters["p2p_gz_ips"] = filters["p2p_gz"]
	return filters
}
