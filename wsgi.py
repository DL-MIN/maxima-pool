"""
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
"""

"""
wsgi.py

@package   MaximaPool
@author    Lars Thoms <lars.thoms@uni-hamburg.de>
@copyright 2020 Universität Hamburg
@license   http://www.gnu.org/copyleft/gpl.html GNU GPL v3 or later
"""

from MaximaPool import app

if __name__ == "__main__":
    app.run()
