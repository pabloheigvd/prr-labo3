# Labo 3

La durée de transmission T est défini := 1.5s

# Test
## Test de non participation a une élection en cours
TODO

## Test seul
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

## Test même élu et départage
On remplace le fichier de configuration `config.json` avec le contenu de `config4departage.test.json`.

On lance quatres processus dans quatre terminaux séparés avec:
```bash
go run process.go <pId>
```
**Note**: on choisit les quatres valeurs de pId allant de 0 à 3

Une fois les quatres processus lancés, on attends 4-5 secondes, puis dans un des terminaux, on tape:
```bash
e
```
pour démarrer une élection. On patiente 3-4 secondes, le temps que l'élection se termine.

On ouvre les fichiers `logs/log<pId>` (`pId` variant de 0 à 3) et consulte la fin de l'output. On vérifie que l'élu le dernier élu est le processus 2. 

On observe dans la configuration utilisée (`configs/config.json`) que le processus 3 aurait pu être élu car il a une aptitude égale à celle du processus 2. Cependant, le processus 3 n'est pas élu car la règle de départage par le plus petit numéro de processus a bien été appliquée. De plus on vérifie que le processus ayant la plus grande aptitude a été élu.

## Test détection de la panne de l'élu
TODO

## Test changement d'aptitude interactif
TODO