# IFLY — Cahier des charges et consignes pour un agent de développement

> **Mission immédiate :** reprendre le projet Wails existant et implémenter un premier prototype vertical (*prototype B*) : importer une image, l'analyser localement, en déduire une composition musicale et faire évoluer réellement l'image et le son au cours d'une expérience temporelle. **Ne commence pas par créer un éditeur graphique complet.**
>
> **Ce document est autonome :** considère-le comme la source de vérité fonctionnelle. Examine le dépôt réel avant de modifier des fichiers et distingue ce qui est déjà implémenté de ce qui est demandé.

## 1. Intention du projet

IFLY est une application **desktop créative et expérimentale**, développée principalement en Go, permettant de **dessiner librement une œuvre**, de la **valider définitivement**, puis de **regarder et d'écouter son évolution irréversible**.

Après validation, le programme :

1. comprend autant que possible le **contenu et l'atmosphère** de l'image ;
2. compose **automatiquement** une musique pertinente, sans intervention musicale de l'utilisateur ;
3. laisse une période de **contemplation**, pendant laquelle l'œuvre reste intacte ;
4. démarre ensuite une **entropie visuelle et musicale** progressive ;
5. termine l'expérience lorsque l'œuvre a suffisamment évolué ou disparu.

Ce n'est **ni un produit commercial ni un séquenceur MIDI**, ni un outil de musique où l'utilisateur assigne un instrument à chaque objet. C'est une expérience artistique : le dessin est libre et la musique est l'interprétation du programme.

### Philosophie non négociable

**« On ne peut pas remonter le temps. »**

- Avant validation : l'utilisateur peut librement dessiner, effacer et corriger.
- **Après validation :** l'œuvre n'est **plus éditable**. Aucun retour au mode dessin, aucun annuler/rétablir, aucun contrôle permettant de retrouver un état passé de cette œuvre.
- L'expérience avance dans une seule direction. **Pas de scrub arrière, de bouton Reverse ou de remise à zéro d'une œuvre validée.**
- L'application peut conserver en interne une **source immuable** pour calculer correctement les transformations. Cette nécessité technique ne doit jamais devenir une commande de restauration accessible à l'utilisateur.
- Les exports sont des **souvenirs** (images, audio, vidéo), pas des projets réouvrables pour modifier l'œuvre validée.
- Une nouvelle création est une **nouvelle œuvre**, jamais une réédition d'une ancienne.
- Une pause temporaire avant reprise est acceptable comme commande de lecture (prévue par le croquis initial) ; elle ne doit jamais permettre de reculer ou de rejouer l'expérience comme si elle n'avait pas eu lieu. La sémantique exacte de la pause et de la reprise après fermeture doit être explicite et cohérente.

Ne prétends pas imposer une « impossibilité cryptographique » de modifier des fichiers sur l'ordinateur de l'utilisateur : il s'agit d'une **règle du comportement de l'application**, à faire respecter durablement et à tester.

## 2. État de départ et environnement

- **OS de développement : Windows 11**.
- **IDE : GoLand**, terminal PowerShell.
- **Projet existant :** `C:\Users\Amogousse\Documents\Dev\ifly`.
- **Wails v3.1.0** a été installé et le projet a été initialisé avec le **template React + TypeScript**.
- Le projet généré a d'abord été placé dans un dossier `ifly/ifly`, puis l'utilisateur a corrigé l'arborescence. **Ne recrée pas de dossier `ifly` imbriqué.** Vérifie la structure effective.
- **Go** est le langage principal côté application/native ; **TypeScript** est acceptable pour l'interface et les API Web de rendu/audio.
- Le projet dispose normalement des fichiers Wails issus du template (par exemple `main.go`, `frontend/`, configuration de build). **Ne suppose pas leur contenu ni leurs chemins exacts : inspecte-les avant tout changement.**

### Stack recommandée (à garder pragmatique)

| Responsabilité | Choix de départ |
|---|---|
| Shell desktop | Wails v3 |
| Backend, gestion de l'œuvre, analyse, orchestration | Go |
| Interface | React + TypeScript |
| Style | CSS personnalisé / CSS Modules ; pas de grosse bibliothèque UI sans besoin démontré |
| Affichage du prototype | Canvas 2D ; WebGL seulement si nécessaire |
| Musique synthétisée | Web Audio API / AudioWorklet si nécessaire |
| Analyse sémantique hors ligne | Petit VLM quantifié, exécuté localement, à benchmarker |
| Distribution finale | Application Windows autonome, modèle et runtime locaux embarqués |

**Évite le sur-engineering :** pas de serveur distant obligatoire, pas de compte, pas de base de données externe, pas d'architecture microservices, pas de Docker, pas de bibliothèque de composants imposée par habitude.

## 3. Ce que doit être l'application finale (vision cible)

### 3.1 Interface visuelle

Le croquis de référence montre une fenêtre desktop avec :

- une **barre d'outils verticale à gauche**, inspirée dans l'esprit des outils de dessin classiques ;
- un **grand espace de dessin central**, prioritaire dans la composition ;
- une zone discrète **actions / export / commandes** en haut à gauche de l'espace central ;
- les commandes **Play / Pause** en haut à droite.

Le design doit rester sobre, compact et adapté au desktop, sans esthétique de dashboard SaaS. L'interface doit valoriser l'image.

### 3.2 Dessin libre (phase ultérieure, pas maintenant)

À terme : pinceau, gomme, couleurs, épaisseur, fond, quelques outils de dessin simples, zoom et déplacement si utiles. L'utilisateur peut dessiner **n'importe quoi** : paysages, portraits, scènes narratives, abstractions, etc.

**Il n'y a PAS de système de calques prévu.** Ne l'ajoute pas. Pas de fonction d'affectation manuelle de sons, instruments, gammes ou notes aux éléments du dessin. Pas de piano-roll, pas de lecteur MIDI comme concept de l'application.

### 3.3 Analyse du sens

La couleur et la géométrie **ne suffisent pas**. Le moteur doit aussi chercher à comprendre la scène, ses sujets, actions, relations, atmosphère et tension émotionnelle.

Exemple : une scène de torture ou de souffrance ne doit pas être classée automatiquement comme joyeuse simplement parce qu'elle contient des couleurs claires. Une opposition volontaire entre violence et douceur peut être pertinente artistiquement, mais elle doit résulter d'une décision de composition cohérente, pas d'une incompréhension triviale.

Ne réduis pas l'interprétation à une règle « couleur rouge = colère » ni à un simple score de sentiment. Les dessins amateurs, simplifiés et abstraits sont difficiles pour les petits modèles : gérer honnêtement l'ambiguïté.

### 3.4 Génération musicale automatique

Le programme choisit **entièrement** la composition : ambiance, harmonie, rythme, notes, timbres, intensité, évolution. L'utilisateur n'a pas à sélectionner ses instruments ou à associer un son à une forme.

La musique doit paraître **composée**, pas être une succession arbitraire de notes aléatoires. Préférer une génération algorithmique contrôlée : motifs, phrases, gamme/centre tonal selon l'interprétation, harmonie, densité, dynamique, tension, règles de développement temporel.

### 3.5 Temps, contemplation et entropie

**L'entropie ne démarre pas immédiatement.** Après analyse et validation, l'image peut être contemplée intacte pendant une durée dépendant, au moins partiellement, de la **complexité globale de l'œuvre**. La musique accompagne déjà cette phase.

La complexité peut combiner : richesse des éléments, densité de contours, diversité visuelle, organisation spatiale et complexité sémantique lorsque celle-ci est fiable. Un dessin minimaliste doit quand même bénéficier d'une durée minimale de contemplation. Prévoir des bornes supérieures et une relation non linéaire pour éviter les durées absurdes.

L'entropie doit **transformer véritablement l'image**, pas seulement superposer une animation de particules à une image initiale qui reste identique jusqu'à la fin : érosion, diffusion, déplacement de régions, désagrégation, perte de structure, décoloration, dissolution, reconfiguration, selon des règles cohérentes.

La musique **évolue en correspondance** avec les transformations visuelles : déstructuration des motifs, perte progressive de stabilité rythmique/harmonique, modulation du timbre, fragmentation, silences, etc. **Un même état d'évolution** pilote son et image. Éviter les effets indépendants sans relation perceptible.

La complexité détermine notamment **quand et combien de temps** l'œuvre évolue ; le contenu et la structure de l'image déterminent **comment** elle évolue.

## 4. Travail immédiat : prototype B, avant l'éditeur

Construire d'abord une chaîne démontrable complète. Le prototype B accepte **une image importée** (PNG, JPEG, WebP si pris en charge), uniquement pour accélérer les tests. Ce n'est pas un engagement d'ajouter l'import au parcours final de dessin.

### Parcours utilisateur du prototype

1. Ouvrir IFLY et **importer une image** locale.
2. Voir l'image en grand, correctement redimensionnée sans déformation.
3. Cliquer sur **Analyser** : analyse visuelle et sémantique, avec progression/état et gestion des erreurs.
4. Voir un **résumé de l'interprétation** pour le diagnostic du prototype (pas nécessairement dans l'interface finale) : description, sujets, atmosphère, tension, métriques visuelles, complexité estimée.
5. Cliquer sur **Valider / Commencer l'expérience**. À partir de là, l'œuvre ne peut plus être modifiée ou réimportée dans **cette session d'œuvre** ; toute nouvelle expérimentation démarre une nouvelle œuvre.
6. Le programme génère et joue une **musique originale synthétisée automatiquement**.
7. Le dessin reste intact durant une **phase de contemplation audible**.
8. L'entropie démarre et transforme progressivement **l'image et la composition**.
9. L'expérience aboutit à une fin perceptible ; l'utilisateur peut alors commencer une **nouvelle œuvre distincte**.

Pour les tests techniques, autoriser un **mode laboratoire séparé** où le développeur peut rejouer la même image / utiliser des paramètres accélérés. Cette possibilité de diagnostic ne doit pas contaminer le comportement de l'expérience finale destinée à l'utilisateur.

### Priorités et livrables incrémentaux

#### B0 — Audit et lancement du projet

- Examiner le dépôt existant avant de modifier quoi que ce soit.
- Vérifier `go.mod`, `frontend/package.json`, les points d'entrée Wails et les commandes réellement disponibles.
- Confirmer la compilation et le lancement avec les outils du projet (`wails3 dev` ou commande officielle appropriée).
- Ne pas réinitialiser le projet ni créer de sous-dossier supplémentaire.

#### B1 — Import et affichage

- Une interface simple, sombre, desktop, avec import, aperçu et zone de diagnostic.
- Contrôle des extensions et du **type réel** des images, gestion des fichiers corrompus, limites raisonnables de taille et de dimensions (éviter les bombes de décompression).
- Gérer correctement rotation/orientation de l'image, proportion, zoom de prévisualisation si nécessaire.
- Libérer les ressources après changement d'image, erreurs et fermeture.

**Terminé lorsque :** une image valide s'affiche ; les entrées non valides produisent un message clair sans planter l'application.

#### B2 — Analyse locale / expérimentation VLM

- Définir un contrat Go **indépendant du moteur d'inférence**, permettant de remplacer un modèle sans réécrire le frontend.
- Commencer avec une analyse **structurelle déterministe** de l'image (palette, contraste, contours, densité, histogrammes, etc.).
- Intégrer un **VLM local** léger pour la description sémantique, sans API payante. **SmolVLM 256M et 500M** sont des candidats à **tester**, pas un choix définitivement validé.
- Étudier un runtime compatible avec Windows, CPU et distribution autonome, par exemple `llama.cpp` et des poids GGUF compatibles ; vérifier réellement la prise en charge multimodale de la variante choisie (encodeur visuel inclus).
- Ne pas supposer qu'un petit VLM comprend correctement des dessins amateurs. Comparer descriptions, erreurs, latence et mémoire sur plusieurs cas : paysage, forme abstraite, personnage simplifié, scène joyeuse, scène menaçante ou violente.
- Préférer une sortie structurée **validée côté Go**. Si le VLM ne respecte pas le schéma ou invente des éléments, prévoir une reprise limitée ou un résultat « incertain » explicite.
- En cas d'indisponibilité de l'IA, montrer une erreur exploitable, éventuellement proposer une analyse structurelle limitée **identifiée comme telle**. Ne pas déguiser des résultats fictifs en résultats d'un modèle.

Exemple de **schéma interne visé**, pas promesse de sortie exacte du modèle :

```json
{
  "description": "A storm approaching a forest",
  "subjects": ["forest", "storm"],
  "mood": ["ominous", "tense"],
  "visualComplexity": 0.65,
  "semanticAmbiguity": "low",
  "musicFeatures": {
    "valence": -0.5,
    "arousal": 0.7,
    "tension": 0.8
  }
}
```

Les champs numériques représentent des **caractéristiques de composition dérivées**, pas des probabilités garanties ou une mesure universelle de l'émotion. Valider et borner ces valeurs.

**Terminé lorsque :** un vrai modèle local est testé sur CPU, ses limites sont documentées dans le compte rendu de l'agent, et son analyse alimente une représentation exploitable par Go.

#### B3 — Premier compositeur musical

- À partir de l'analyse, générer une **structure musicale reproductible** avec une graine liée à l'œuvre.
- Établir des règles pour tempo, gamme ou centre tonal, progression/motifs, densité, timbres synthétiques, dynamique et tension.
- Synthétiser localement via Web Audio API. Pas d'assets protégés, pas d'API musicale externe.
- Démarrer la musique après validation, avec vraie prise en compte des contraintes de démarrage du contexte audio Web (interaction utilisateur).
- Éviter les clics audio, les notes bloquées, les accumulations d'oscillateurs et les boucles qui dérivent.

**Terminé lorsque :** des images de caractères visuels/sémantiques différents produisent des pièces perceptiblement différentes, audibles et non chaotiques par défaut.

#### B4 — Entropie visuelle réelle et synchronisation

- Utiliser une **source image immuable** et un **état temporel monotone** pour piloter l'évolution.
- Phase de contemplation **sans dégradation**, puis évolution visuelle significative, puis fin.
- Première version : transformations procédurales de l'image rendue, basées sur ses caractéristiques (par exemple segmentation simple couleur/luminance, blocs déplacés, érosion de contours, dissolution, diffusion). La segmentation sémantique précise pourra venir plus tard ; **un VLM décrivant l'image ne fournit pas automatiquement des masques par objet**.
- Faire évoluer les paramètres musicaux avec les **mêmes variables temporelles / événements d'entropie** que le rendu.
- Éviter que le Canvas conserve constamment en arrière-plan l'image intacte sous un effet cosmétique. La dégradation doit modifier le contenu visuel affiché de manière observable.
- Lors d'une pause, arrêter l'avancement de l'horloge de simulation et du moteur musical sans autoriser de retour arrière. À la reprise, continuer à partir de l'état courant.

**Terminé lorsque :** une même image passe de contemplation à désorganisation visuelle et musicale, avec transitions audibles et visibles synchronisées.

#### B5 — Validation, persistance et robustesse

- Machine à états explicite, par exemple :

```text
DRAFT -> ANALYZING -> READY -> CONTEMPLATING -> DECAYING -> ENDED
```

  `DRAFT` existe pour l'import de test ; la version finale aura un état de dessin libre avant validation.
- **Pas de transition vers `DRAFT` pour la même œuvre après `CONTEMPLATING`.**
- Refuser les actions incompatibles avec l'état courant, y compris dans les services Go, pas seulement en désactivant des boutons.
- Conserver minimalement l'identifiant, la graine, le statut, les métriques, les horodatages et l'évolution nécessaire pour récupérer proprement une session.
- Définir explicitement la reprise après fermeture : **ne jamais remonter une œuvre à un état antérieur**. Pour le premier prototype, une politique simple de reprise de l'état sauvegardé sans rejouer le passé est préférable à une simulation de temps réel écoulé en arrière-plan non maîtrisée. Documenter ce choix dans le compte rendu ; aucun bouton de réinitialisation pour l'œuvre validée.

### Exports (après stabilisation du prototype)

- Priorité : **PNG** d'un instant courant.
- Ensuite : **WAV** de la composition.
- Plus tard : **vidéo** synchronisée, si le coût et la complexité sont acceptables.
- L'export n'est pas un moyen de réimporter l'ancienne œuvre pour la modifier.
- Aucun système de fichiers de projet éditable après validation n'est attendu.

## 5. Architecture indicative, sans rigidité excessive

Garder les responsabilités séparées, sans multiplier les packages à l'avance :

```text
React/TypeScript UI
  - import, preview, diagnostics, experience controls
  - canvas renderer (not React state per pixel/frame)
  - Web Audio synthesis and scheduling
          |
          | Wails bindings/events
          v
Go application
  - validation and lifecycle state machine
  - immutable artwork metadata
  - structural image analysis
  - semantic analysis abstraction
  - musical plan / entropy rules and shared parameters
          |
          v
Local VLM inference process/runtime
  - model files bundled for offline use
```

**Performance :** ne pas envoyer chaque frame Canvas ou échantillon audio à travers la frontière Go/JavaScript. Les calculs temps réel doivent rester dans leur moteur de rendu/audio. Go peut produire une analyse initiale, un plan musical, des paramètres d'évolution et des événements à faible fréquence.

**Concurrence :** analyse IA asynchrone, avec annulation, délai maximal, limitation du nombre d'analyses simultanées, et fermeture propre du processus enfant. L'interface ne doit pas se bloquer pendant l'inférence CPU.

**Distribution autonome :** à terme, inclure l'exécutable d'inférence, le modèle quantifié et l'encodeur visuel nécessaires, avec vérification des licences et notices requises. Aucune installation d'Ollama, de Python, de Docker ni connexion à un service cloud ne doit être exigée pour utiliser l'application finale. Pour le développement, un runtime local externe peut être utilisé temporairement **si cela est clairement signalé**, puis remplacé par un packaging autonome.

**Sécurité :** traiter images et sorties du modèle comme des données non fiables. Vérifier types, limites mémoire, chemin de lecture, taille du contenu et JSON retourné. Ne jamais interpoler du texte utilisateur dans une commande shell. Si un processus local expose temporairement un endpoint, le lier uniquement à l'interface loopback et ne pas le rendre accessible sur le réseau.

## 6. Tests et critères de réussite

- L'application démarre avec Wails sur Windows 11.
- L'import de plusieurs images différentes fonctionne ; images invalides/énormes correctement refusées.
- Un **vrai modèle** local produit des analyses exploitables sans appel Internet ; les résultats factices sont isolés aux tests.
- Sur un corpus de croquis rudimentaires, relever les erreurs de reconnaissance et les cas incertains. Ne pas inventer de taux de réussite.
- Une image menaçante/violente et un paysage paisible ne doivent pas produire exactement la même interprétation et le même plan musical par défaut.
- La musique est générée automatiquement, sans configuration de notes ou d'instruments par l'utilisateur.
- L'image reste **inchangée pendant la contemplation**, puis subit une transformation **réelle**, pas un simple overlay.
- La musique évolue avec les transformations visuelles.
- Après validation, aucune édition de l'œuvre ni retour dans le passé n'est possible dans le parcours normal.
- Une pause/reprise ne réinitialise pas l'horloge.
- Les erreurs d'analyse ou de synthèse ne font pas planter l'application.
- Vérifier les fuites de ressources (timers, AudioNodes, images, contexte Web, processus d'inférence).

## 7. Règles de travail pour l'agent

1. **Inspecter le projet existant avant de coder.** Ne pas inventer de fichiers, chemins ou API Wails. Respecter la version réellement installée et la structure effective.
2. **Implémenter directement**, pas seulement proposer une architecture théorique. Commencer par **B0 puis B1** et poursuivre avec B2 dès que la base fonctionne.
3. Faire des changements **petits, cohérents, vérifiables** ; compiler et tester dès que possible. Ne pas remplacer aveuglément les fichiers du template sans les avoir lus.
4. **Ne pas ajouter de calques, séquenceur MIDI, choix manuel d'instrument, marketplace, authentification, backend cloud ou design system lourd.**
5. Code, noms, commentaires et logs **en anglais**. Interface utilisateur et explications destinées au créateur **en français**.
6. Garder un code lisible, idiomatique, modulaire seulement lorsque nécessaire ; gérer erreurs et annulations ; éviter les dépendances non justifiées.
7. **Ne pas utiliser Git**, ne pas créer de `.gitignore`, `README` ni autre document annexe. Respecter la structure présente, sans effacer les fichiers de l'utilisateur.
8. Si tu fournis des fichiers modifiés au lieu de travailler directement dans le projet : **1 à 3 fichiers = fichiers individuels** ; **4 fichiers ou plus = une archive ZIP contenant uniquement ce qui doit être ajouté/modifié**, en conservant leurs chemins relatifs.
9. À chaque jalon, rendre compte de : fichiers réellement modifiés, fonctionnement vérifié, limites connues, commandes exactes de lancement et prochaine priorité. **Ne prétends pas avoir exécuté un test que tu n'as pas exécuté.**
10. En cas de doute technique, consulter la documentation officielle correspondant aux versions installées. Pour l'IA locale, mesurer la qualité sur de vrais dessins avant de figer le modèle.

## 8. Hors périmètre immédiat

- Éditeur complet de dessin : **plus tard** (prototype A).
- Calques : **non prévus**.
- IA cloud / API facturée : **non**.
- Génération complète d'images par diffusion : **non** ; l'utilisateur crée l'œuvre, le programme la transforme.
- Segmentation sémantique précise d'objets : **éventuelle évolution**, si nécessaire et réalisable localement.
- Fonction de retour arrière, reverse entropy, modification après validation : **interdites**.
- Système d'édition MIDI, mapping manuel des sons, contrôles d'instrument utilisateur : **interdits comme concept central**.
- Export vidéo, installer/packager la distribution finale, design sophistiqué : **après validation du cœur expérimental**.

## 9. Première instruction à exécuter maintenant

> Ouvre le dépôt existant d'IFLY. Identifie précisément les fichiers Wails/Go/React présents et confirme comment lancer l'application. Sans réinitialiser le projet, réalise le premier jalon du prototype B : **import d'une image locale et affichage fiable dans une interface desktop minimale**. Prépare la séparation entre analyse visuelle Go, futur VLM local et frontend. Ensuite, établis un test réel et mesurable des petits VLM sur CPU avant de développer le compositeur musical ou l'éditeur de dessin. Respecte toutes les contraintes d'irréversibilité, mais n'impose pas encore la validation définitive sur l'outil d'import de laboratoire tant que le pipeline n'est pas fonctionnel.
