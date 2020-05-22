#
# MaximaPool.py
#
# @author Lars Thoms <lars.thoms@uni-hamburg.de>
#

from flask import Flask, request, Response
import subprocess
import os
import io
import zipfile
import tempfile

# Default route of flask application
app = Flask(__name__)


@app.route('/', methods=['POST'])
def maxima_request():
    instructions = {

        # Maxima commands
        'input': request.form.get('input'),

        # Timeout in seconds (10s as default max)
        'timeout': min(int(request.form.get('timeout')), 10000) / 1000,

        # Unused parameters sent by qtype_stack
        'ploturlbase': request.form.get('ploturlbase'),
        'version': request.form.get('version')
    }

    # Reject empty requests
    if instructions['input'] is None:
        return bad_request()

    try:

        # Spawn Maxima process
        with subprocess.Popen(['maxima', '--quiet'], stdin=subprocess.PIPE, stdout=subprocess.PIPE) as process:

            # Create temporary directory for plots
            with tempfile.TemporaryDirectory(dir='plots/') as plot_directory:
                try:

                    # Replace plot path in STACK configurations and add Maxima commands
                    maxima_input = stack_config.replace('%PLOT-DIR%', plot_directory + '/') + instructions['input']

                    # Start calculating with timeout
                    result = process.communicate(
                        input=maxima_input.encode(),
                        timeout=instructions['timeout'])

                # Terminate process if expired
                except subprocess.TimeoutExpired:
                    process.kill()
                    return bad_request()

                # Successful run of Maxima
                if process.returncode is 0:

                    # No plots? Return result as plaintext
                    if not os.listdir(plot_directory):
                        return Response(response=result[0].decode(), status=200, mimetype='text/plain',
                                        content_type='text/plain')

                    # Build zip due to plots
                    else:

                        # Initialize zip buffer to avoid unnecessary disk usage
                        zip_buffer = io.BytesIO()
                        with zipfile.ZipFile(zip_buffer, 'a', zipfile.ZIP_DEFLATED, False) as zip_file:

                            # Dump Maxima result to zip archive
                            zip_file.writestr('OUTPUT', result[0])

                            # Add plots to zip archive
                            for root, dirs, files in os.walk(plot_directory):
                                for filename in files:
                                    zip_file.write(os.path.join(root, filename), arcname=filename)

                        # Response with zip file
                        return Response(response=zip_buffer.getvalue(), status=200, mimetype='application/zip',
                                        content_type='application/zip')

                # Unsuccessful run of Maxima
                else:
                    return bad_request()

    # Something went really wrong?
    except subprocess.CalledProcessError:
        return bad_request()


# Return HTTP 416 as recognized as failed run in qtype_stack
def bad_request():
    return Response(response='👎', status=416)


# Load STACK configuration
def get_stack_config():
    with open('assets/maximalocal.mac') as config_file:
        return config_file.read()


# Cache of STACK configuration
stack_config = get_stack_config()
