
Contexte Personnel : Je suis architecte informatique et dev. Je suis ouvert à toute proposition d’amélioration pour avoir une architecture modulaire, performante et maintenable.

Contexte technique : Tout doit être très rapide et le code doit être en mono-repo. 

Contexte Application: Yousoon est une application de sortie qui propose aux utilisateurs des réductions à chaque sortie aux utilisateurs. Les fournisseurs (bars, restaurants, organismes de sorties) proposent leurs réductions à un groupe.  Les utilisateurs choisissent leurs sorties et s’enregistrent pour la sortie. 
Le but de Yousoon est de proposer des sorties aux utilisateurs avec des promotions.
1) J'apporte des clients à mes partenaires qui en échange feront des réductions aux personnes que je leur apporterai
2) J'offre à mes clients des sorties avec des réductions, via mes partenaires et en échange il me paient via un abonnement
3) Je fais l'intermédiaire en allant chercher des clients et des partenaires et en les mettant en relation. Mes clients veulent sortir pour pas trop cher, et mes partenaires veulent des clients


Ce que je veux :
- un site web en react typescript pour que les fournisseurs puissent proposer leurs offres. Il sera déployé sur Kubernetes. Il communique en GraphQL et le BE est en Go (même que celui de l’App mobile)
- Une App mobile en Flutter avec le design Figma de Yousoon (NomFichier) avec les contraintes suivantes :
    - Tu dois respecter scrupuleusement le design. 
    - L’App doit communiquer avec un backend en Go via le protocole GraphQL  que je pourrais déployer sur Kubernetes 
    - Le Backend doit avoir une architecture micros services et doit répondre rapidement. Il faut donc du cache (coté App Flutter, coté BDD, coté services). Les requêtes doivent répondre en moins de 50ms. Je veux pouvoir tester la latence des requêtes (quel outil ?) via des tests de charge 
    - Je veux des composants flutter dès qu’on peut mutualiser
    - L’App doit utiliser un service de validation des carte d’identité
    - Je veux des tests unitaires et e2e pour l’inscription (je pourrais en ajouter d’autres après)
- La BDD du site Fournisseur et de l’App cliente   est en MongoDB et sera à déployer sur Kubernetes. L’App fournisseur partage ses données avec l’App Flutter. Les partenaires font leur offres et l’App Yousoon les propose aux clients. 
- Je veux un site web avec une section Fournisseurs et une prestation de l’App. Il sera déployé sur Kubernetes 
- Les offres partenaires sont par types de partenaire car ils auront des offres différentes. 

Avant de commencer, pose moi les questions que tu aurai pour affiner ou améliorer la stock technique. Et aussi sur le fonctionnel des l’applications. Ou toute autre question sur le projet. Dis moi ce que j’aurai pu oublier. Je veux des questions techniques, donc n’hésite pas à rentrer dans le détail si nécessaire. Liste moi les points critiques à adresser. 

Ne fais pas de code pour le moment, on améliore le besoin et spécifications. 

Reformule ce prompt avant de poser tes questions afin qu’il soit bien structuré 
Si tu peux prendre un prompte depuis un fichier alors met le en fichier où il faut et dis moi où c’est ou comment faire. 

Le site web serra de présentation sera www.yousoon.com et la section partenaires sera business.yousoon.com 

Je veux une app la plus native et réactive que possible. 
Je veux que tu détailles bien dans le prompt
 la stack technique et et ce que tu comptes faire. 

Je veux que tu dises quel model de données tu comptes faire en fonction du design (App et Site partenaires). Mets ça dans un fichier et tu pourra les lier au fichier du prompt. 

Après : pour le design, il y a des écrans ou c’est des sliders. Je veux que tu me disent lesquels tu penses que c’est. Si tu ne comprends pas bien un écran, poses tes questions. Liste moi les composants réutilisables que tu penses faire 

Si tu peux avoir des prompt via fichier alors je veux que les prompts soient par modules/dossiers (site vitrine, app, site partenaires). 
Si tu veux tu peux dans ces dossiers faires un fichier avec le prompt dédié et un fichier avec les questions pour améliorer. N’oublie pas que je veux tu technique et pertinent. 


# Code génération

En respectant les recommandations d'architecture validées.
En respectant le design system et les écrans 

Crée moi un plan avec les étapes permettant de générer : 
    - l'Application Mobile :
        - Les microservices backend
        - l'application mobile
    - Le site d'admin
    - Le site partenaire

Crée un plan de génération par module (App, backend, partenaires, site vitrine, etc..) dans chaque module.


# Génération Full

Avec ce que tu m'as écrit 
"Génère le Backend Phase 1 : Infrastructure commune en suivant le plan backend/GENERATION_PLAN.md"

Fais moi les prompts pour te demande de tout générer.
Tu dois mettre à jour les status à chaque fois et valider que tu n'as pas déjà fait une étape. Le but est que si la génération se plante , que tu puisse repartir de ou tu étais 