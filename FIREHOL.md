# Index differences from FireHOL

In general, I corrected various URLs from their index. Some URLs are still using HTTP, but I will try to check if they offer an HTTPS endpoint, given some free time.

Some services are still in the list with a note about their status (e.g., Gateway Timeout) because they may come back online in the near future.

## Added

These lists have been added.

* urlvir_last
* urlhaus
* acknowledged_scanners (network scanners, not malicious, but still)

## Discontinued

These lists have been removed because they have been discontinued, or their website does not exist anymore.

* badipscom
* normshield
* urandomusto
* lashback_ubl
* ipblacklistcloud
* all hpHosts (hosts-file.net)
* hphosts_hjk
* rosinstrument.com
* proxyrss.com
* proxylists.net
* torstatus.blutmagie.de
* taichung
* snort_ipfilter (moved to Talos Intel)
* nt_ssh_7d
* nt_malware_irc
* nt_malware_http
* nt_malware_dns
* gofferje_sip
* dshield_top_1000
* zeus_badips
* zeus
* all Abuse.ch under ransomwaretracker.abuse.ch (moved to urlhaus)
* malwaredomainlist
* malc0de
* all PacketMail.net
* nullsecure
* antispam.imp.ch

## Went commercial

These list went commercial, meaning that now some sort of subscription is needed, or they are free only to some categories of users. They have been removed, but maybe one day I may re-add them and find a way to include the description of the license in the YAML.

* bambenek
* CoinBlockerLists

MaxMind it's free behind authentication. However, they have their downloaded. Maybe we can still implement that (provide some way to pass a token?).

* geolite2_asn
* geolite2_country

## Old lists

* eSentire "malfeed" is not updated since 2016. Is it really useful?
    * https://github.com/eSentire/malfeed
