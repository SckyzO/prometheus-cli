import requests
from tabulate import tabulate
import readline

# URL de l'API Prometheus
prometheus_url = "http://localhost:9090/api/v1"

# Fonction pour récupérer les métriques disponibles
def get_metrics():
    url = f"{prometheus_url}/label/__name__/values"
    response = requests.get(url)
    data = response.json()
    return data['data']

# Fonction pour effectuer une requête Prometheus
def query_prometheus(query):
    url = f"{prometheus_url}/query"
    params = {
        'query': query
    }
    response = requests.get(url, params=params)
    data = response.json()
    return data['data']['result']

# Fonction pour afficher les résultats sous forme de tableau
def display_table(results):
    headers = ['Métrique'] + sorted(set().union(*[result['metric'].keys() for result in results]))
    headers.remove('__name__')  # Supprime la deuxième colonne '__name__'
    headers.append('Valeur')  # Ajoute la colonne 'Valeur'
    table = []

    for result in results:
        row = [result['metric'].get(label, '') for label in headers[1:-1]]
        value = result['value'][1]  # Récupère la valeur
        row.append(value)  # Ajoute la valeur en dernière colonne
        row.insert(0, result['metric']['__name__'])
        table.append(row)

    print(tabulate(table, headers, tablefmt="grid"))

# Récupération des métriques disponibles
metrics = get_metrics()

# Utilisation de readline pour l'autocomplétion
readline.set_completer_delims('\t')
readline.parse_and_bind("tab: complete")

# Ajout des métriques à l'autocomplétion
readline.set_completer(lambda text, state: [metric for metric in metrics if metric.startswith(text)][state])

# Demande à l'utilisateur de saisir une requête
query = input("Entrez une requête Prometheus : ")

# Exécution de la requête
results = query_prometheus(query)

# Affichage des résultats dans un tableau
display_table(results)
