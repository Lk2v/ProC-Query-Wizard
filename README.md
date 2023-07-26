# ProC-Query-Wizard

Ce mini-programme Golang a été développé en une demi-journée dans le cadre de mon stage pour automatiser la transformation de requêtes Pro*C en requêtes SQL testables. Il est doté d'une interface graphique simple comportant deux champs de texte.

Le premier champ de texte permet de saisir la requête Pro*C en entrée, et le second champ de texte affiche le résultat de la transformation. La transformation de la requête consiste en plusieurs étapes :

1. Remplacement des noms de variables par leurs valeurs, qui sont extraites à partir du fichier commun.h.
2. Désactivation des clauses INTO des requêtes SELECT en les mettant en commentaire.
3. Remplacement des occurrences de structName.fields par :field.

## Fonctionnalités

- Saisie de la requête Pro*C dans l'interface.
- Transformation automatique de la requête pour la rendre testable via SQL.
- Remplacement des noms de variables par leurs valeurs correspondantes à partir du fichier commun.h.
- Désactivation des clauses INTO pour les requêtes SELECT en les mettant en commentaire.
- Remplacement des champs de struct par le nom de leur champ.

## Demo

![Screenshot](/Documentation/Screenshots/1.png)

> Interface du programme

### Exemple

Voici un exemple pour la requête source suivante :
```SQL
EXEC SQL SELECT PER_IDT,
                TYPE_INTERVENTION,
                SOC_PER_IDT
         INTO   :l_ts_Sortie->ms_nb_adh_as,
                :l_ts_Sortie->ms_nb_adh_at,
                :l_ts_Sortie->ms_nb_adh_av
         FROM   DENOMBRE_RISQUES_SOCIETE DEN
         WHERE  :l_ts_Entree->me_typ_per = :g_TYP_PER_S
         AND    DEN.SOC_PER_IDT = :l_ts_Entree->me_per_idt
         AND TYPE_INTERVENTION = :g_INTERV_RAD 
         AND STATUS_INTERVALE = :g_STATUS_INTERVALE
```
Pour un fichier commun.h avec les valeurs completé le resultat sera :

```SQL
EXEC SQL SELECT PER_IDT,
                TYPE_INTERVENTION,
                SOC_PER_IDT
         --INTO   :ms_nb_adh_as,
         --       :ms_nb_adh_at,
         --       :ms_nb_adh_av
         FROM   DENOMBRE_RISQUES_SOCIETE DEN
         WHERE  :me_typ_per = "T_ADHERENT_01"
         AND    DEN.SOC_PER_IDT = :me_per_idt
         AND TYPE_INTERVENTION = "CODE_682"
         AND STATUS_INTERVALE = "DNB_RSQ_SOC_INTERVALE"
```

## Installation

1. Assurez-vous d'avoir Golang installé sur votre système.
2. Clonez ce dépôt GitHub
3. Accédez au répertoire du projet
4. Exécutez le programme : go run .

## Configuration

Aucune configuration particulière n'est nécessaire pour utiliser ce programme. Cependant, assurez-vous que le fichier commun.h contenant les valeurs des variables est présent dans le même répertoire que le programme.