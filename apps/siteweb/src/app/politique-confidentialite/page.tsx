import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Politique de Confidentialité',
  description: 'Politique de confidentialité et protection des données personnelles Yousoon',
};

export default function PolitiqueConfidentialitePage() {
  return (
    <div className="container mx-auto px-4 py-24">
      <div className="mx-auto max-w-3xl">
        <h1 className="text-4xl font-bold">Politique de Confidentialité</h1>
        <p className="mt-4 text-muted-foreground">
          Dernière mise à jour : Décembre 2024
        </p>

        <div className="mt-12 space-y-8">
          <section>
            <h2 className="text-2xl font-semibold">1. Introduction</h2>
            <p className="mt-4 text-muted-foreground">
              Yousoon SAS (&quot;nous&quot;, &quot;notre&quot;, &quot;Yousoon&quot;) s&apos;engage à protéger la vie
              privée de ses utilisateurs. Cette politique de confidentialité
              explique comment nous collectons, utilisons, partageons et
              protégeons vos informations personnelles.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">
              2. Données collectées
            </h2>
            <p className="mt-4 text-muted-foreground">
              Nous collectons les données suivantes :
            </p>
            <ul className="mt-4 list-disc space-y-2 pl-6 text-muted-foreground">
              <li>
                <strong className="text-foreground">
                  Données d&apos;identification :
                </strong>{' '}
                nom, prénom, email, numéro de téléphone
              </li>
              <li>
                <strong className="text-foreground">
                  Données de vérification d&apos;identité :
                </strong>{' '}
                copie de pièce d&apos;identité (CNI, passeport)
              </li>
              <li>
                <strong className="text-foreground">
                  Données de géolocalisation :
                </strong>{' '}
                position pour afficher les offres à proximité
              </li>
              <li>
                <strong className="text-foreground">
                  Données d&apos;utilisation :
                </strong>{' '}
                historique des réservations, favoris, préférences
              </li>
              <li>
                <strong className="text-foreground">
                  Données techniques :
                </strong>{' '}
                adresse IP, type d&apos;appareil, système d&apos;exploitation
              </li>
            </ul>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">
              3. Finalités du traitement
            </h2>
            <p className="mt-4 text-muted-foreground">
              Vos données sont utilisées pour :
            </p>
            <ul className="mt-4 list-disc space-y-2 pl-6 text-muted-foreground">
              <li>Créer et gérer votre compte utilisateur</li>
              <li>Vérifier votre identité (obligation légale)</li>
              <li>Vous proposer des offres personnalisées</li>
              <li>Traiter vos réservations</li>
              <li>Vous envoyer des notifications (avec votre consentement)</li>
              <li>Améliorer nos services</li>
              <li>Assurer la sécurité de la plateforme</li>
            </ul>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">4. Base légale</h2>
            <p className="mt-4 text-muted-foreground">
              Nous traitons vos données sur les bases légales suivantes :
            </p>
            <ul className="mt-4 list-disc space-y-2 pl-6 text-muted-foreground">
              <li>
                <strong className="text-foreground">Contrat :</strong>{' '}
                exécution du service auquel vous avez souscrit
              </li>
              <li>
                <strong className="text-foreground">Consentement :</strong>{' '}
                notifications marketing, géolocalisation
              </li>
              <li>
                <strong className="text-foreground">
                  Obligation légale :
                </strong>{' '}
                vérification d&apos;identité
              </li>
              <li>
                <strong className="text-foreground">
                  Intérêt légitime :
                </strong>{' '}
                amélioration des services, sécurité
              </li>
            </ul>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">
              5. Durée de conservation
            </h2>
            <p className="mt-4 text-muted-foreground">
              Vos données sont conservées pendant la durée nécessaire aux
              finalités pour lesquelles elles ont été collectées :
            </p>
            <ul className="mt-4 list-disc space-y-2 pl-6 text-muted-foreground">
              <li>Données de compte : durée de vie du compte + 3 ans</li>
              <li>Documents d&apos;identité : 5 ans après vérification</li>
              <li>Données de transaction : 10 ans (obligation légale)</li>
              <li>Logs techniques : 1 an</li>
            </ul>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">6. Partage des données</h2>
            <p className="mt-4 text-muted-foreground">
              Nous partageons vos données uniquement avec :
            </p>
            <ul className="mt-4 list-disc space-y-2 pl-6 text-muted-foreground">
              <li>
                <strong className="text-foreground">Partenaires :</strong>{' '}
                informations nécessaires à la réservation
              </li>
              <li>
                <strong className="text-foreground">
                  Sous-traitants techniques :
                </strong>{' '}
                hébergement (AWS), analytics (Amplitude)
              </li>
              <li>
                <strong className="text-foreground">Autorités :</strong> sur
                demande légale
              </li>
            </ul>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">7. Vos droits</h2>
            <p className="mt-4 text-muted-foreground">
              Conformément au RGPD, vous disposez des droits suivants :
            </p>
            <ul className="mt-4 list-disc space-y-2 pl-6 text-muted-foreground">
              <li>
                <strong className="text-foreground">Droit d&apos;accès :</strong>{' '}
                obtenir une copie de vos données
              </li>
              <li>
                <strong className="text-foreground">
                  Droit de rectification :
                </strong>{' '}
                corriger vos données
              </li>
              <li>
                <strong className="text-foreground">
                  Droit à l&apos;effacement :
                </strong>{' '}
                supprimer vos données
              </li>
              <li>
                <strong className="text-foreground">
                  Droit à la portabilité :
                </strong>{' '}
                récupérer vos données
              </li>
              <li>
                <strong className="text-foreground">
                  Droit d&apos;opposition :
                </strong>{' '}
                vous opposer au traitement
              </li>
              <li>
                <strong className="text-foreground">
                  Droit de limitation :
                </strong>{' '}
                limiter le traitement
              </li>
            </ul>
            <p className="mt-4 text-muted-foreground">
              Pour exercer vos droits, contactez-nous à{' '}
              <a
                href="mailto:dpo@yousoon.com"
                className="text-primary hover:underline"
              >
                dpo@yousoon.com
              </a>
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">
              8. Suppression de compte (RGPD)
            </h2>
            <p className="mt-4 text-muted-foreground">
              Vous pouvez demander la suppression de votre compte à tout moment
              depuis l&apos;application ou par email. Une période de grâce de{' '}
              <strong className="text-foreground">30 jours</strong> vous permet
              d&apos;annuler cette demande et récupérer votre compte. Passé ce
              délai, toutes vos données personnelles seront supprimées de
              manière irréversible.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">9. Sécurité</h2>
            <p className="mt-4 text-muted-foreground">
              Nous mettons en œuvre des mesures de sécurité appropriées pour
              protéger vos données : chiffrement SSL/TLS, authentification
              sécurisée, audits réguliers, hébergement en Europe (conformité
              RGPD).
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">10. Cookies</h2>
            <p className="mt-4 text-muted-foreground">
              Nous utilisons des cookies essentiels au fonctionnement du site et
              des cookies analytics (avec votre consentement) pour améliorer nos
              services.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">11. Contact</h2>
            <p className="mt-4 text-muted-foreground">
              Pour toute question relative à cette politique, contactez notre
              Délégué à la Protection des Données :{' '}
              <a
                href="mailto:dpo@yousoon.com"
                className="text-primary hover:underline"
              >
                dpo@yousoon.com
              </a>
            </p>
            <p className="mt-4 text-muted-foreground">
              Vous pouvez également déposer une réclamation auprès de la CNIL :{' '}
              <a
                href="https://www.cnil.fr"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:underline"
              >
                www.cnil.fr
              </a>
            </p>
          </section>
        </div>
      </div>
    </div>
  );
}
