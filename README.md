# reportAbuse

For further information about this project, please visit the blog article at http://blog.inspirant.de/index.php?controller=post&action=view&id_post=49.

# Configuration

First make a copy of the `config_template.json` and name it `config.json` or any other name you like. Now adapt the settings and enter credentials.

# Calling the tool

There are comandline parameters to use:

`reportAbuse -config=<file> { -debug }`

**-config** to set the configuration file to use. Use like -config=config.json.

**-debug** to stop the tool from remembering the IPs and not sending to the hosters abuse address (just send to what is defined in **smtpCopy**).

## Configuration options

**logfiles** is an array of logfiles to parse. They must be all of the same format.

**mode** is either "page" or "direct". The page mode means that a typical scan page must have been scanned by the attacker (see detection.go for the scanned page names). In direct mode, all RegExp hits are seen as an attack (eg if your log already only lists attacks).

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

## Examples

### Scan some apache log for typical pages

Scan file ssl_error_log_short and use given RegEx to get values. Use **page** mode (attacked page name must be one of the list) and report only if there were at least two tries to attack the website. A single try does not trigger the email.

    { 
        "logfiles": [ 
            "ssl_error_log_short"
        ],

        "databasePath": "abuse_log.db",

        "myPage": "my.page.com",

        "regExp": "(?m)\\[(.*?)\\].* \\[client (.*):\\d*\\].*'.*\/opt\/provider\/REGIFY_PUBLIC(.*)'",
        "regGroupDate": 1,
        "regGroupIP": 2,
        "regGroupPage": 3,

        "mode": "page",

        "minAttacks": 2,
        "smtpHost": "sslout.df.eu:25",
        "smtpUser": "sender@df.eu",
        "smtpPwd": "password",
        "smtpCopy": "sender@df.eu"
    }

### Scan a log whwre RegExp only catches known attacks directly (IDS)

Scan file apache_error_log and use given RegEx to get values. Use **direct** mode (every RegEx hit means an attack) and report only if there were at least two tries to attack the website. A single try does not trigger the email.

    { 
        "logfiles": [ 
            "apache_error_log"
        ],

        "databasePath": "abuse_log.db",

        "myPage": "my.page.com",

        "regExp": "(?msU)DATE: (\\d\\d\\d\\d-\\d\\d-\\d\\d\\s\\d\\d:\\d\\d:\\d\\d).*CALLER:\\s(\\d+\\.\\d+\\.\\d+\\.\\d+).*URL:\\s(.*)(?:\\?|$)",
        "regGroupDate": 1,
        "regGroupIP": 2,
        "regGroupPage": 3,

        "mode": "direct",

        "minAttacks": 2,
        "smtpHost": "sslout.df.eu:25",
        "smtpUser": "sender@df.eu",
        "smtpPwd": "password",
        "smtpCopy": "sender@df.eu"
    }
