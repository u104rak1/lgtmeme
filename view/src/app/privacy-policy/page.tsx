import { TWITTER_LINK_ENDPOINT } from "@/utils/constants";
import { css } from "@@/styled-system/css";

const PrivacyPolicy = () => {
  return (
    <div className={containerCss}>
      <div className={contentCss}>
        <h1 className={h1Css}>Privacy Policy</h1>
        <br />
        <p>
          {
            'This Privacy Policy (hereinafter referred to as "Policy") describes how the service provider (hereinafter referred to as "Operator") of the LGTM image generation service "LGTMeme" (hereinafter referred to as "the Service") collects, uses, and protects the personal information of users (hereinafter referred to as "Users"). Before using the Service, please carefully read and agree to the following policy.'
          }
        </p>
        <h2 className={h2Css}>1. Information Collected</h2>
        <p>1.1. Information provided by Users</p>
        <p>{"・User's name (optional)"}</p>
        <p>{"・User's email address (optional)"}</p>
        <p>1.2. Automatically collected information</p>
        <p>{"・User's IP address"}</p>
        <p>・Browser information</p>
        <p>{"・User's device information"}</p>
        <h2 className={h2Css}>2. Purpose of Collection</h2>
        <p>
          2.1. The Operator uses the collected information for the following
          purposes:
        </p>
        <p>・To provide and improve the Service</p>
        <p>・To provide support to Users</p>
        <p>・To collect and analyze anonymous statistical data</p>
        <p>
          ・To ensure the security of the Service and compliance with legal
          requirements
        </p>
        <h2 className={h2Css}>3. Sharing of Information</h2>
        <p>
          {
            "3.1. The Operator does not share Users' personal information with third parties."
          }
        </p>
        <h2 className={h2Css}>4. Cookies and Tracking Technologies</h2>
        <p>
          {
            "4.1. Since the Service does not use cookies or similar tracking technologies, no information related to Users' privacy is collected."
          }
        </p>
        <h2 className={h2Css}>5. User Choices</h2>
        <p>
          5.1. Users can use the Service anonymously. Information provided (such
          as name and email address) is voluntary. Users have the choice whether
          to provide this information or not.
        </p>
        <h2 className={h2Css}>6. Data Retention Period</h2>
        <p>
          6.1. The Operator retains the collected information for the necessary
          period and then deletes or anonymizes the information once that period
          has ended.
        </p>
        <h2 className={h2Css}>7. Privacy of Minors</h2>
        <p>
          7.1. The Service is not provided to minors and does not intentionally
          collect information from minors.
        </p>
        <h2 className={h2Css}>8. Changes to the Privacy Policy</h2>
        <p>
          8.1. The Operator reserves the right to change the Policy at any time.
          There is no obligation to notify Users of changes. Users are
          responsible for regularly reviewing the Policy.
        </p>
        <h2 className={h2Css}>9. Contact Information</h2>
        <p>
          9.1. For questions, requests, and other communications regarding the
          Privacy Policy, please contact the following:
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

export default PrivacyPolicy;
