# reportAbuse

For further information about this project, please visit the blog article at http://blog.inspirant.de/index.php?controller=post&action=view&id_post=49.

# Configuration

First make a copy of the `config_template.json` and name it `config.json`. Now adapt the settings and enter credentials.

## Configuration options

**logfiles** is an array of logfiles to parse. They must be all of the same format.

**databasePath** sets the path and name of the database folder to use for storing reported IPs. This is important because without, you would send repeated and multiple reports to the same hoster or provider all the time. Make sure you have write permissions!

**myPage** is simply the name of your website. It is used in the message to the hoster.

**regExp** is the escaped regular expression (regex) to find three things in your logfile: Date&Time, IP address and the page causing the error. 

My log entry looks like this:
`[Tue Feb 28 11:25:58.518511 2023] [:error] [pid 51189] [client 142.93.243.6:54338] script '/opt/provider/REGIFY_PUBLIC/wso112233.php' not found or unable to stat, referer: http://free.regify.com/wso112233.php`

Because of this, I did my regex like this (without json escaping):

`(?m)\[(.*?)\].* \[client (.*):\d*\].*'.*\/opt\/provider\/REGIFY_PUBLIC(.*)'`

It will use the golang modifier for _global_ and _multiline_ followed by the date in square brackets and the other values I'm interested in. I strongly recommend to use a tool like https://regex101.com/ for making your regex!

Dont forget that you need to escape backslash in json.

**regGroupDate** sets the result number of your capturing group for the date/time component.

**regGroupIP** sets the result number of your capturing group for the IP address.

**regGroupPage** sets the result number of your capturing group for the called page. Please make sure you also get the leading slash. So your group returns `/wp-login.php` including the slash.

**minAttacks** defines, how many tries the attacker need to do until this script generates an abuse email. You can set this to 1 if you like to report every try, even single ones. But I prefer to report only three or more. This reduces false positives and only reports bigger scans.

**smtpHost** is used to set your smtp server including port for sending the abuse emails. Set like "smtp.mycompany.com:25".

**smtpUser** is the login email address and also the shown sender.

**smtpPwd** is the login password.

**smtpCopy** is an email address you like to get copies to (eg your own). Mainly for controlling the progress and mails.