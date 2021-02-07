# Laboratoire 3

La durée de transmission T est défini := 1.5s
# Test
## Test seul
Un test manuel est de lancer un processus et de voir si le message annonçant le 
résultat de l'élection intial correspond à quoi on s'attend.

On choisit pId := 0

On lance la commande:

```bash
go run process.go <pId>
```
Après 3-4 secondes d'attente, on apercoit:
```bash
L'elu de l'election initiale est le processus: 0
```
On peut répéter le test avec des numéros de processus différents et s'apercevoir 
que le numéro du processus gagnant l'élection change en conséquence.

## Test de non participation a une élection en cours
On remplace le fichier de configuration `config.json` avec le contenu de `config4departage.test.json`.

On lance un processus avec `<pId> := 0`. Quand le texte "*Lancement d'une nouvelle election*" 
apparait, alors on lance un autre processus avec `<piD> := 1`. 
On observe que dans le premier terminal, le processus 0 est élu lors du premier cycle d'élection. 
Peu de temps après, le deuxième processus indique qu'il a été élu. 

On constate que le premier processus voit le deuxième processus élu en tapant dans 
la console du premier processus:
```bash
g
```

## Test même élu et départage
On remplace le fichier de configuration `config.json` avec le contenu de `config4departage.test.json`.

On lance quatres processus dans quatre terminaux séparés avec:
```bash
go run process.go <pId>
```
**Note**: on choisit les quatres valeurs de pId allant de 0 à 3

Une fois les quatres processus lancés, on attends 4-5 secondes que l'élection se termine.

On ouvre les fichiers `logs/log<pId>` (`pId` variant de 0 à 3) et consulte la fin 
de l'output. On vérifie que l'élu le dernier élu est le processus 2. 
On peut aussi taper `g` dans chacun des terminaux pour le vérifier.

On observe dans la configuration utilisée (`configs/config.json`) que le 
processus 3 aurait pu être élu car il a une aptitude égale à celle du processus 2. 
Cependant, le processus 3 n'est pas élu car la règle de départage par le plus 
petit numéro de processus a bien été appliquée. 
De plus on vérifie que le processus ayant la plus grande aptitude a été élu.

## Test détection de la panne de l'élu
On remplace le fichier de configuration `config.json` avec le contenu de 
`config4departage.test.json`.

On lance quatres processus dans quatre terminaux séparés avec:
```bash
go run process.go <pId>
```
**Note**: on choisit les quatres valeurs de pId allant de 0 à 3

Une fois les quatres processus lancés, on attends 4-5 secondes.

Une fois le premier cycle d'élection terminé, on va fermer le processus 2 avec `ctrl-c`.

Après 3-4 secondes, on vérifie que le processus élu est le processus 3 en tapant 
```bash
g
```
dans les trois terminaux restants.

On répète l'opération en fermant le processus 3, vérifie que le processus 1 
est élu avec `g`.

On vérifie que le processus 0 est élu après avoir fermé le processus `1`.

## Test changement d'aptitude interactif
On remplace le fichier de configuration `config.json` avec le contenu de `config4departage.test.json`.

On lance quatres processus dans quatre terminaux séparés avec:
```bash
go run process.go <pId>
```
**Note**: on choisit les quatres valeurs de pId allant de 0 à 3

Une fois les quatres processus lancés, on attends 4-5 secondes, le temps que l'élection 
se termine.

On prend un terminal au hasard et on tape une aptitude grande, e.g. 123.

On patiente 3-4 secondes, le temps que l'aptitude se propage aux autres processus.

On vérifie que le processus pour lequel l'aptitude est 123 est élu en tapant 
dans chaque terminaux:
```bash
g
```

On modifie l'aptitude a 1 et on vérifie que ce même processus n'est plus 
élu de manière similaire.

## Test un changement d'aptitude pas avant que l'élection en cours se termine
On remplace le fichier de configuration `config.json` avec le contenu de 
`config4departage.test.json`.

On lance quatres processus dans quatre terminaux séparés avec :
```bash
go run process.go <pId>
```
**Note**: on choisit les quatre valeurs de pId allant de 0 à 3

Une fois les quatre processus lancés, on attend 4-5 secondes, le temps que l'élection 
se termine.

On modifie l'aptitude d'un terminal avec 123. 1 seconde après avoir modifié 
l'aptitude, on remodifie l'aptitude avec 1234.

On aperçoit dans la console que deux élections ont été lancées.