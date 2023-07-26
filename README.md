# SQL-Query-Wizard

Ce programme Golang a été développé dans le cadre de mon stage pour automatiser la transformation de requêtes Pro*C en requêtes SQL testables. Il est doté d'une interface graphique simple comportant deux champs de texte.

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

## Installation

1. Assurez-vous d'avoir Golang installé sur votre système.
2. Clonez ce dépôt GitHub
3. Accédez au répertoire du projet
4. Exécutez le programme : go run .

## Configuration

Aucune configuration particulière n'est nécessaire pour utiliser ce programme. Cependant, assurez-vous que le fichier commun.h contenant les valeurs des variables est présent dans le même répertoire que le programme.
