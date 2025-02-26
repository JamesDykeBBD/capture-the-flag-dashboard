from flask import Flask, request, render_template_string, render_template
import os

app = Flask(__name__)

# Create a secret flag
FLAG = os.environ.get('FLAG', 'CTF{test_flag_please_change_in_docker}')

# Define flag location in user's home directory
FLAG_PATH = '/home/ctfuser/flag.txt'

# Create directory if it doesn't exist
os.makedirs(os.path.dirname(FLAG_PATH), exist_ok=True)

# Write flag to a file
with open(FLAG_PATH, 'w') as f:
    f.write(FLAG)

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/result')
def result():
    # Get the name parameter from the request
    name = request.args.get('name', '')
    
    # VULNERABLE CODE: directly inserting user input into the template
    template = '''
    <!DOCTYPE html>
    <html>
    <head>
        <title>CTF Challenge - SSTI in Jinja2</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                max-width: 800px;
                margin: 0 auto;
                padding: 20px;
                background-color: #f5f5f5;
            }
            .container {
                background-color: white;
                padding: 20px;
                border-radius: 8px;
                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            }
            h1 {
                color: #333;
            }
            .result {
                margin-top: 20px;
                padding: 10px;
                background-color: #e8f4fc;
                border-left: 4px solid #4da6ff;
            }
            .footer {
                margin-top: 20px;
                text-align: center;
                font-size: 0.8em;
                color: #666;
            }
            a {
                color: #4CAF50;
                text-decoration: none;
            }
            a:hover {
                text-decoration: underline;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Result</h1>
            
            <div class="result">
                <h3>Hello, ''' + name + '''!</h3>
                <p>Thanks for visiting our site.</p>
            </div>
            
            <p><a href="/">Try another name</a></p>
        </div>
        
        <div class="footer">
            <p>CTF Challenge - Server Side Template Injection in Jinja2</p>
        </div>
    </body>
    </html>
    '''
    
    # Render the template with the user's input
    return render_template_string(template)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=False)
