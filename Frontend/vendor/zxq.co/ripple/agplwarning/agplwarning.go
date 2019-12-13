package agplwarning

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"
)

var pages = [...]string{
	`This is a friendly reminder that the GNU AGPL adds an additional clause to
the standard GNU GPL, which is that you MUST distribute the source code for the
software once you publish it on the web.
    This is not to be considered professional legal advice. For further
information, refer to the LICENSE file which contains the whole license, or ask
your lawyer. If you did not receive a copy of the LICENSE file with this
software, you can refer to the online version:
    https://www.gnu.org/licenses/agpl-3.0.html`,
	`In order to comply with the license, should you have made any modification
to the original copy of the software, which should contain a link to the
source code, however minor it is, you are under the legal obligation to provide
the source code once you publish the software on the Web.
    Another obligation is that of stating your changes. This is usually done by
cloning the original git repository of the project and stating your changes
through the creation of commits, which allow us to determine when a specific
change was done.`,
	`Furthermore, all the original clauses of the GNU General Public License
are kept intact, which means you have the obligation to
    * Keep the AGPL License, without possibility of sublicensing the software
      or making it available under any other more liberal license.
    * Keep the copyright notice of the original authors
    Failure to do so will result in a request to follow the License, and
repeated violation of the license could result in a legal fight.`,
	`For more information on the FSF and software freedom, refer to:
    * What is free software? https://www.gnu.org/philosophy/free-sw.html
    * Free Software Is Even More Important Now
      https://www.gnu.org/philosophy/free-software-even-more-important.html
    * The GNU operating system https://www.gnu.org
    * The Free Software Foundation https://www.fsf.org
    Thank you for reading this and following our license terms.`,
}

// Warn shows a warning about the GNU Affero General Public License the first
// time the software is run. The state is saved in
// ~/.config/[namespace]_license_agreed.
func Warn(namespace, projectName string) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("initialization of agplwarning failed: %v", err)
	}
	if err := os.MkdirAll(usr.HomeDir+"/.config", 0755); err != nil {
		return fmt.Errorf("can't create config dir: %v", err)
	}
	agreedFilename := usr.HomeDir + "/.config/" + namespace + "_license_agreed"
	if _, err := os.Stat(agreedFilename); !os.IsNotExist(err) {
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	// file does not exist. Show warning.
	fmt.Printf("    %s, and most/all software related to %s,\n"+
		"is licensed under the GNU Affero General Public License.\n\n", projectName, namespace)
	for _, page := range pages {
		fmt.Println("    " + page)
		fmt.Println("\nPress Enter to continue")
		_, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
	}
	fmt.Println("Please write 'I agree' to accept the terms of the license.")
	res, err := reader.ReadString('\n')
	if err != nil || !strings.Contains(strings.ToLower(res), "i agree") {
		fmt.Println("License not agreed. Quitting.")
		os.Exit(1)
	}
	f, err := os.Create(agreedFilename)
	if err != nil {
		return fmt.Errorf("couldn't save read status: %v", err)
	}
	return f.Close()
}
