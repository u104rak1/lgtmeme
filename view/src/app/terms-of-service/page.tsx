import { PAGE_ENDPOINTS, TWITTER_LINK_ENDPOINT } from "@/utils/constants";
import { css } from "@@/styled-system/css";

const TermsOfService = () => {
  return (
    <div className={containerCss}>
      <div className={contentCss}>
        <h1 className={h1Css}>Terms of Service</h1>
        <br />
        <p>
          {
            'These Terms of Service (hereinafter referred to as "Terms") establish the conditions for the use of the LGTM image generation service "LGTMeme" (hereinafter referred to as "the Service") provided by the service provider (hereinafter referred to as "Operator") between the Operator and the service user (hereinafter referred to as "User"). Before using the Service, you must read and agree to the following terms.'
          }
        </p>
        <h2 className={h2Css}>1. Provision of Service</h2>
        <p>
          1.1. The Operator provides the Service but reserves the right to
          change the content and available features of the Service at any time.
          The Operator may decide to change or discontinue the Service without
          notice.
        </p>
        <p>
          1.2. Users shall have all rights necessary for providing the Service.
          Users must ensure that they do not infringe on the intellectual
          property rights of others, such as copyright, trademark, or patent
          rights, when using the Service.
        </p>
        <h2 className={h2Css}>2. Disclaimer</h2>
        <p>
          {
            '2.1. The Service is provided "as is". The Operator makes no warranty as to the accuracy, completeness, or applicability of the Service. Users assume all risk associated with the use of the Service, and the Operator is not liable for any damages.'
          }
        </p>
        <h2 className={h2Css}>3. User Responsibilities</h2>
        <p>
          3.1. When using the Service, Users agree to and shall follow these
          conditions:
        </p>
        <p>・Not to misuse the Service and engage in illegal activities.</p>
        <p>
          ・To pay attention to and not violate the copyright or license of the
          source images.
        </p>
        <p>・Not to send excessive requests that could overload the Service.</p>
        <p>
          {
            "・Not to create images deemed inappropriate by the Operator. Inappropriate images may be removed at the Operator's discretion."
          }
        </p>
        <h2 className={h2Css}>4. Privacy</h2>
        <p>
          {
            "4.1. User privacy will be in accordance with the Operator's privacy policy. Please refer to the privacy policy "
          }
          <a className={linkCss} href={PAGE_ENDPOINTS.privacyPolicy}>
            here
          </a>
          .
        </p>
        <h2 className={h2Css}>5. Legal Dispute Resolution</h2>
        <p>
          5.1. Disputes related to these Terms shall be governed by Japanese
          law, and the Tokyo District Court shall be the exclusive court of
          first instance.
        </p>
        <h2 className={h2Css}>6. Modification of Terms</h2>
        <p>
          6.1. The Operator reserves the right to modify these Terms at any
          time. There is no obligation to notify Users of changes. Users are
          responsible for regularly reviewing these Terms.
        </p>
        <h2 className={h2Css}>7. Contact Information</h2>
        <p>
          7.1. For inquiries, reports, and other communications regarding the
          Service, please contact the following:
        </p>
        <br />
        <a className={linkCss} href={TWITTER_LINK_ENDPOINT} target="_blank">
          {TWITTER_LINK_ENDPOINT}
        </a>
      </div>
    </div>
  );
};

const containerCss = css({
  display: "flex",
  justifyContent: "center",
  padding: "5",
});
const contentCss = css({ maxWidth: "1090px", padding: 5 });
const h1Css = css({ fontSize: "2xl", fontWeight: "bold" });
const h2Css = css({ fontSize: "xl", fontWeight: "bold", marginTop: "5" });
const linkCss = css({
  color: "LIGHT_BLUE",
  _hover: {
    opacity: 0.8,
    borderBottom: "1px solid",
    borderColor: "LIGHT_BLUE",
  },
});

export default TermsOfService;
