import os
from pathlib import Path

PAGES = (
    """This is a friendly reminder that the GNU AGPL adds an additional clause to
the standard GNU GPL, which is that you MUST distribute the source code for the
software once you publish it on the web.
    This is not to be considered professional legal advice. For further
information, refer to the LICENSE file which contains the whole license, or ask
your lawyer. If you did not receive a copy of the LICENSE file with this
software, you can refer to the online version:
    https://www.gnu.org/licenses/agpl-3.0.html""",
    """In order to comply with the license, should you have made any modification
to the original copy of the software, which should contain a link to the
source code, however minor it is, you are under the legal obligation to provide
the source code once you publish the software on the Web.
    Another obligation is that of stating your changes. This is usually done by
cloning the original git repository of the project and stating your changes
through the creation of commits, which allow us to determine when a specific
change was done.""",
    """Furthermore, all the original clauses of the GNU General Public License
are kept intact, which means you have the obligation to
    * Keep the AGPL License, without possibility of sublicensing the software
      or making it available under any other more liberal license.
    * Keep the copyright notice of the original authors
    Failure to do so will result in a request to follow the License, and
repeated violation of the license could result in a legal fight.""",
    """For more information on the FSF and software freedom, refer to:
    * What is free software? https://www.gnu.org/philosophy/free-sw.html
    * Free Software Is Even More Important Now
      https://www.gnu.org/philosophy/free-software-even-more-important.html
    * The GNU operating system https://www.gnu.org
    * The Free Software Foundation https://www.fsf.org
    Thank you for reading this and following our license terms.""",
)


class LicenseError(Exception):
    pass


def check_license(namespace, project_name):
    license_folder_path = "{}/.config/".format(Path.home())
    if not os.path.isdir(license_folder_path):
        try:
            os.mkdir(license_folder_path, mode=0o755)
        except OSError as e:
            raise LicenseError("Cannot create .config dir: {}".format(e))
    agreed_file_name = "{}/{}_license_agreed".format(license_folder_path, namespace)
    if os.path.isfile(agreed_file_name):
        return

    print(
        "    {}, and most/all software related to {},\n"
        "is licensed under the GNU Affero General Public License.\n".format(project_name, namespace)
    )
    for page in PAGES:
        print(" " * 4 + "\n" + page)
        try:
            input("\nPress Enter to continue")
        except KeyboardInterrupt:
            raise LicenseError("License not read. Quitting.")

    if input("\nPlease write 'I agree' to accept the terms of the license.\n").lower().strip() != "i agree":
        raise LicenseError("License not agreed. Quitting.")

    try:
        open(agreed_file_name, "a").close()
    except IOError as e:
        raise LicenseError("Couldn't save read status: {}".format(e))
