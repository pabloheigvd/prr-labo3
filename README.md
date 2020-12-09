# prr-labo3

# Test seul
Un test manuel est de lancer un processus et de voir voir si le message annonçant le résultat de l'élection intial correspond à quoi on s'attend.

On choisit pId := 0

On lance la commande:

```bash
go run process.go <pId>
```
Après 3-4 secondes d'attente, on apercoit:
```bash
L'elu de l'election initiale est le processus: 0
```
On peut répéter le test avec des numéros de processus différents et s'apercevoir que le numéro du processus gagnant l'élection change en conséquence.