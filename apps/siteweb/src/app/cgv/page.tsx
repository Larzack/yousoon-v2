import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Conditions Générales de Vente',
  description: 'Conditions générales de vente et d\'utilisation de Yousoon',
};

export default function CGVPage() {
  return (
    <div className="container mx-auto px-4 py-24">
      <div className="mx-auto max-w-3xl">
        <h1 className="text-4xl font-bold">Conditions Générales de Vente</h1>
        <p className="mt-4 text-muted-foreground">
          Dernière mise à jour : Décembre 2024
        </p>

        <div className="mt-12 space-y-8">
          <section>
            <h2 className="text-2xl font-semibold">1. Objet</h2>
            <p className="mt-4 text-muted-foreground">
              Les présentes Conditions Générales de Vente (CGV) régissent les
              relations contractuelles entre Yousoon SAS et ses utilisateurs
              pour l&apos;utilisation de l&apos;application mobile Yousoon et des
              services associés.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">2. Services proposés</h2>
            <p className="mt-4 text-muted-foreground">
              Yousoon est une plateforme de sorties avec réductions qui met en
              relation des utilisateurs (&quot;Yousooners&quot;) avec des établissements
              partenaires (bars, restaurants, lieux de loisirs) proposant des
              offres à prix réduits.
            </p>
            <p className="mt-4 text-muted-foreground">
              Les services comprennent :
            </p>
            <ul className="mt-4 list-disc space-y-2 pl-6 text-muted-foreground">
              <li>Accès au catalogue d&apos;offres géolocalisées</li>
              <li>Réservation d&apos;offres</li>
              <li>Validation par QR code chez les partenaires</li>
              <li>Gestion des favoris et préférences</li>
              <li>Notifications personnalisées</li>
            </ul>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">3. Abonnements</h2>
            <h3 className="mt-6 text-xl font-medium">3.1 Types d&apos;abonnements</h3>
            <ul className="mt-4 list-disc space-y-2 pl-6 text-muted-foreground">
              <li>
                <strong className="text-foreground">Gratuit :</strong> accès
                limité aux offres (3 réservations/mois)
              </li>
              <li>
                <strong className="text-foreground">Premium mensuel :</strong>{' '}
                accès illimité aux offres
              </li>
              <li>
                <strong className="text-foreground">Premium annuel :</strong>{' '}
                accès illimité + 2 mois offerts
              </li>
            </ul>

            <h3 className="mt-6 text-xl font-medium">3.2 Période d&apos;essai</h3>
            <p className="mt-4 text-muted-foreground">
              Une période d&apos;essai gratuite de 30 jours est proposée pour les
              nouveaux utilisateurs. À l&apos;issue de cette période, l&apos;abonnement
              payant débute automatiquement sauf annulation préalable.
            </p>

            <h3 className="mt-6 text-xl font-medium">3.3 Paiement</h3>
            <p className="mt-4 text-muted-foreground">
              Tous les paiements sont effectués via Apple Pay (iOS) ou Google
              Play (Android). Yousoon n&apos;a pas accès à vos informations de
              paiement.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">
              4. Conditions d&apos;utilisation
            </h2>
            <h3 className="mt-6 text-xl font-medium">4.1 Éligibilité</h3>
            <p className="mt-4 text-muted-foreground">
              Pour utiliser Yousoon, vous devez :
            </p>
            <ul className="mt-4 list-disc space-y-2 pl-6 text-muted-foreground">
              <li>Être âgé d&apos;au moins 18 ans</li>
              <li>Disposer d&apos;une pièce d&apos;identité valide</li>
              <li>Compléter la vérification d&apos;identité</li>
            </ul>

            <h3 className="mt-6 text-xl font-medium">
              4.2 Vérification d&apos;identité
            </h3>
            <p className="mt-4 text-muted-foreground">
              La vérification d&apos;identité est obligatoire pour effectuer des
              réservations. Elle s&apos;effectue par scan de pièce d&apos;identité (CNI,
              passeport ou permis de conduire). Vous disposez de 10 tentatives
              maximum.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">5. Réservations</h2>
            <h3 className="mt-6 text-xl font-medium">5.1 Processus</h3>
            <p className="mt-4 text-muted-foreground">
              Une réservation est valide pendant 30 minutes après sa création.
              Passé ce délai, elle expire automatiquement.
            </p>

            <h3 className="mt-6 text-xl font-medium">5.2 Check-in</h3>
            <p className="mt-4 text-muted-foreground">
              Le check-in s&apos;effectue exclusivement par présentation du QR code
              généré dans l&apos;application chez le partenaire. Un QR code ne peut
              être utilisé qu&apos;une seule fois.
            </p>

            <h3 className="mt-6 text-xl font-medium">5.3 Annulation</h3>
            <p className="mt-4 text-muted-foreground">
              Vous pouvez annuler une réservation à tout moment avant le
              check-in. Une fois le check-in effectué, aucune annulation n&apos;est
              possible.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">6. Droit de rétractation</h2>
            <p className="mt-4 text-muted-foreground">
              Conformément à l&apos;article L221-28 du Code de la consommation, le
              droit de rétractation ne s&apos;applique pas aux services pleinement
              exécutés avant la fin du délai de rétractation avec l&apos;accord du
              consommateur.
            </p>
            <p className="mt-4 text-muted-foreground">
              Pour l&apos;abonnement Premium, vous pouvez annuler à tout moment depuis
              les paramètres de votre store (App Store ou Google Play).
              L&apos;abonnement reste actif jusqu&apos;à la fin de la période payée.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">
              7. Responsabilités
            </h2>
            <h3 className="mt-6 text-xl font-medium">7.1 Yousoon</h3>
            <p className="mt-4 text-muted-foreground">
              Yousoon agit en tant qu&apos;intermédiaire entre les utilisateurs et
              les partenaires. Nous ne sommes pas responsables de la qualité des
              prestations fournies par les partenaires.
            </p>

            <h3 className="mt-6 text-xl font-medium">7.2 Utilisateur</h3>
            <p className="mt-4 text-muted-foreground">
              L&apos;utilisateur s&apos;engage à utiliser le service de manière loyale et
              conforme aux présentes CGV. Tout abus (réservations multiples non
              utilisées, faux documents, etc.) pourra entraîner la suspension du
              compte.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">8. Propriété intellectuelle</h2>
            <p className="mt-4 text-muted-foreground">
              L&apos;ensemble des éléments de l&apos;application (design, marques, logos,
              textes, code) sont la propriété exclusive de Yousoon SAS. Toute
              reproduction est interdite sans autorisation.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">9. Modification des CGV</h2>
            <p className="mt-4 text-muted-foreground">
              Yousoon se réserve le droit de modifier les présentes CGV. Les
              utilisateurs seront informés par notification dans l&apos;application.
              La poursuite de l&apos;utilisation du service vaut acceptation des
              nouvelles CGV.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">10. Litiges</h2>
            <p className="mt-4 text-muted-foreground">
              En cas de litige, une solution amiable sera recherchée avant toute
              action judiciaire. À défaut, les tribunaux français seront seuls
              compétents.
            </p>
            <p className="mt-4 text-muted-foreground">
              Vous pouvez également recourir à la médiation de la consommation :{' '}
              <a
                href="https://www.economie.gouv.fr/mediation-conso"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:underline"
              >
                www.economie.gouv.fr/mediation-conso
              </a>
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">11. Contact</h2>
            <p className="mt-4 text-muted-foreground">
              Pour toute question relative aux présentes CGV :{' '}
              <a
                href="mailto:contact@yousoon.com"
                className="text-primary hover:underline"
              >
                contact@yousoon.com
              </a>
            </p>
          </section>
        </div>
      </div>
    </div>
  );
}
