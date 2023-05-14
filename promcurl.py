import requests
from tabulate import tabulate
import readline

# Prometheus API URL
prometheus_url = "http://localhost:9090/api/v1"

# Function to get available metrics
def get_metrics():
    url = f"{prometheus_url}/label/__name__/values"
    response = requests.get(url)
    data = response.json()
    return data['data']

# Function to query Prometheus
def query_prometheus(query):
    url = f"{prometheus_url}/query"
    params = {
        'query': query
    }
    response = requests.get(url, params=params)
    data = response.json()
    return data['data']['result']

# Function to display results in a table
def display_table(results):
    headers = ['Metric'] + sorted(set().union(*[result['metric'].keys() for result in results]))
    headers.remove('__name__')  # Remove the duplicate '__name__' column
    headers.append('Value')  # Add the 'Value' column
    table = []

    for result in results:
        row = [result['metric'].get(label, '') for label in headers[1:-1]]
        value = result['value'][1]  # Get the value
        row.append(value)  # Add the value in the last column
        row.insert(0, result['metric']['__name__'])
        table.append(row)

    print(tabulate(table, headers, tablefmt="grid"))

# Get available metrics
metrics = get_metrics()

# Use readline for autocompletion
readline.set_completer_delims('\t')
readline.parse_and_bind("tab: complete")

# Add metrics to autocompletion
readline.set_completer(lambda text, state: [metric for metric in metrics if metric.startswith(text)][state])

# Ask the user to enter a Prometheus query
query = input("Enter a Prometheus query: ")

# Execute the query
results = query_prometheus(query)

# Display the results in a table
display_table(results)
