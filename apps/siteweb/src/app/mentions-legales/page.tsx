import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Mentions Légales',
  description: 'Mentions légales du site Yousoon',
};

export default function MentionsLegalesPage() {
  return (
    <div className="container mx-auto px-4 py-24">
      <div className="mx-auto max-w-3xl">
        <h1 className="text-4xl font-bold">Mentions Légales</h1>
        <p className="mt-4 text-muted-foreground">
          Dernière mise à jour : Décembre 2024
        </p>

        <div className="mt-12 space-y-8">
          <section>
            <h2 className="text-2xl font-semibold">1. Éditeur du site</h2>
            <div className="mt-4 space-y-2 text-muted-foreground">
              <p>
                <strong className="text-foreground">Raison sociale :</strong>{' '}
                Yousoon SAS
              </p>
              <p>
                <strong className="text-foreground">
                  Forme juridique :
                </strong>{' '}
                Société par Actions Simplifiée (SAS)
              </p>
              <p>
                <strong className="text-foreground">
                  Capital social :
                </strong>{' '}
                10 000 €
              </p>
              <p>
                <strong className="text-foreground">Siège social :</strong>{' '}
                Paris, France
              </p>
              <p>
                <strong className="text-foreground">
                  Numéro SIRET :
                </strong>{' '}
                XXX XXX XXX XXXXX
              </p>
              <p>
                <strong className="text-foreground">Numéro TVA :</strong> FR XX
                XXX XXX XXX
              </p>
              <p>
                <strong className="text-foreground">
                  Directeur de la publication :
                </strong>{' '}
                [Nom du directeur]
              </p>
              <p>
                <strong className="text-foreground">Email :</strong>{' '}
                <a
                  href="mailto:contact@yousoon.com"
                  className="text-primary hover:underline"
                >
                  contact@yousoon.com
                </a>
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">2. Hébergement</h2>
            <div className="mt-4 space-y-2 text-muted-foreground">
              <p>
                <strong className="text-foreground">Hébergeur :</strong> Amazon
                Web Services (AWS)
              </p>
              <p>
                <strong className="text-foreground">Adresse :</strong> AWS
                EMEA SARL, 38 Avenue John F. Kennedy, L-1855 Luxembourg
              </p>
              <p>
                <strong className="text-foreground">
                  Région d&apos;hébergement :
                </strong>{' '}
                Europe (Irlande) - Conformité RGPD
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">
              3. Propriété intellectuelle
            </h2>
            <p className="mt-4 text-muted-foreground">
              L&apos;ensemble du contenu de ce site (textes, images, logos, graphismes,
              vidéos, etc.) est la propriété exclusive de Yousoon SAS ou de ses
              partenaires et est protégé par les lois françaises et
              internationales relatives à la propriété intellectuelle.
            </p>
            <p className="mt-4 text-muted-foreground">
              Toute reproduction, représentation, modification, publication ou
              adaptation de tout ou partie des éléments du site, quel que soit
              le moyen ou le procédé utilisé, est interdite sans autorisation
              écrite préalable de Yousoon SAS.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">
              4. Données personnelles et RGPD
            </h2>
            <p className="mt-4 text-muted-foreground">
              Conformément au Règlement Général sur la Protection des Données
              (RGPD) et à la loi Informatique et Libertés, vous disposez de
              droits sur vos données personnelles. Pour plus d&apos;informations,
              consultez notre{' '}
              <a
                href="/politique-confidentialite"
                className="text-primary hover:underline"
              >
                Politique de confidentialité
              </a>
              .
            </p>
            <p className="mt-4 text-muted-foreground">
              <strong className="text-foreground">
                Délégué à la Protection des Données (DPO) :
              </strong>{' '}
              <a
                href="mailto:dpo@yousoon.com"
                className="text-primary hover:underline"
              >
                dpo@yousoon.com
              </a>
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">5. Cookies</h2>
            <p className="mt-4 text-muted-foreground">
              Le site utilise des cookies pour améliorer l&apos;expérience
              utilisateur. Vous pouvez configurer votre navigateur pour refuser
              les cookies ou être alerté lorsque des cookies sont envoyés.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">
              6. Limitation de responsabilité
            </h2>
            <p className="mt-4 text-muted-foreground">
              Yousoon SAS s&apos;efforce d&apos;assurer l&apos;exactitude des informations
              diffusées sur ce site, mais ne saurait être tenue pour
              responsable des erreurs, omissions ou résultats qui pourraient
              être obtenus par un mauvais usage de ces informations.
            </p>
          </section>

          <section>
            <h2 className="text-2xl font-semibold">7. Droit applicable</h2>
            <p className="mt-4 text-muted-foreground">
              Les présentes mentions légales sont régies par le droit français.
              En cas de litige, les tribunaux français seront seuls compétents.
            </p>
          </section>
        </div>
      </div>
    </div>
  );
}
