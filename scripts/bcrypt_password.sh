# -b takes the password from the second command argument
# -n prints the hash to stdout instead of writing it to a file
# -B instructs to use bcrypt
# -C 10 sets the bcrypt cost to 10
htpasswd -bnBC 10 $1 $2
