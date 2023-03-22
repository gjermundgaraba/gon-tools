package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/gjermundgaraba/gon/chains"
	"strings"
	"unicode/utf8"
)

type quizQuestion struct {
	question       string
	encryptedPrice string
	flow           string
	class          chains.NFTClass
	nft            chains.NFT
	owner          string
}

func (q quizQuestion) Label() string {
	return fmt.Sprintf("%s %s", q.nft.ID, q.question)
}

var fullQuiz = []quizQuestion{
	{
		question:       "How many official incentivized testnet events did Cosmos held before the Game of NFTs?",
		encryptedPrice: "iJ2L8TzccjgMDiaXDhgPqEV2c_LcguHkaphwfx7mR-yAoelAf-AXI-zX575JAjohu_jkfdX86Q15HdLazVareJfQ4w0-zN4h42OiBdeeCEYRJEttmxsaU7ppzsdNPm1hNaa8CGufyk45gPJYWBFTIToLcK6V1ghuh4cymbc5BvUeipogNpQ6fK0mDrAqrR7evrhMCLNDvDoOlkDMJcVV6ZlIVi5dU4s5A5QQNPUP",
	},
	{
		question:       "How many rounds does the Game of NFTs phase 1 have in total?",
		encryptedPrice: "TJ5dyseXmb4HYXXcCUpWxVaUngByc5S-iyIVfINZ5t6Sl4J5JgxVxCAqRpzV7VsdmjeoYZpukLWQrG_qKccTMmSFOEpQ311IDmPwbd24cvyhqlHXKKu90OiPOAOInYcSHofzbRlFlBU7sSWeA2XlcGz5t1Nk8CUhz4pnCVQ7DCqjd2qErCpganvjqrjcZyoYyWNrx8XlCOcsp-E3cpD6QdMK2dikIh3Ahvw=",
	},
	{
		question:       "How many stages does the Game of NFTs phase 1 have in total?",
		encryptedPrice: "CEShaUvEZG68KzeCrSCW_Gd2yBF3S1VTnkcZ4Kc7fyjEFCOQ6MSzroFoVCdud4kjOLmiFSnSJZ8uqu9wrkYrU6VYLRedJi-4lPWBF2gBKIdI_lKmGiXegM_2DJ6MnpBQG9kPCd-0GU0mhZ8QCYZ2ySwwSPYjj_dnXoEeGUd9OmrmnJGM1ofjfy1BSWOU3oiWQ65s2hbqHmQtNoOa53awBrT1opQ9sWUn5HY=",
	},
	{
		question:       "How many chains are in the Game of NFTs phase 1 incentivized testing?",
		encryptedPrice: "vU7S5MipfGP67mbn9byzRhCU-FJkzjb0PekvbvZCGkl1BaVkjeH__7P1C20zKPMcv3xMgFYUg8yanxFTGQRD7EG-7hP2laB7dSAYVYrWjZ152oJS1G7_QwTWyAsj9Ye1OyLgEvma69Rb7B9aYhgE_eh9xaafpHKqAMfAUeOVvDtUOCV_sZ6V7G8rUx0t3LZdDifvxQhSwHQpBied5p0aUaJauO88hiCQWCE=",
	},
	{
		question:       "How many task flows are in the GoN event?",
		encryptedPrice: "KEN-PD5Wonw8ldbZ22taCyu7neNN4THhUq8k-4zagqaUh8xcVZj1ICNhTNXmz_XJ7v-i7Q-bPhBRbcj9e-wfBVT5azlRRzIqWQuOs67Pdd4YewkugfQwESqF-C_Iiq_ImH4sNJRULqWQHz6LdOrGWDNyDOdgVYsYaOcJHsMIthxpFPCzExvQ7R43CP-q8ZfWi2TBYj0uQlhbI4BR3TQir_V2IjsStfHXl0C1wKFc",
	},
	{
		question:       "How many official channel \u0026 port pairs are in the GoN event?",
		encryptedPrice: "g18kagRONTeJFx_O95r3qAgcF9DT0tzfy7opgfuITKsmmmOQj_xvKcwpQt6JkuwXnfxgUvbOj0cIos3M5iSFFh3gUN7s5ly_hIEDCthXpuu3EbtybNue33cNw7a9zhb56bSPYl4ScvtxIPCX--gzuQd1yMDeh7os6uPjUf0hPervWr2AK0Mb4mAQz-LRXKjklD-7oEOR_nI60Qno0temk1dOXghQAN4=",
	},
	{
		question:       "What's the chain id of IRISnet in GoN testnests?",
		encryptedPrice: "UR306yZofYQGiCtGdfnwdxYhOLTD9C7H2It1qA7DUwh57Ijg5BWt2szFb-twoEsdkghqwMZO88w6MNdxQzeGPrPfmz6hpa5wlxasvFafkN2j3JfdEA_zHIVkMORP1v_ALeiP9W-ymYj6M5FcubW8JIO4EZrXNRaCvm8sTbCI_FWIm9zkl4qWQib4pPjmhbO4TT7unkgV1ZlcblPDy4PTkQyUeWKOKp9iKkE=",
	},
	{
		question:       "What's the chain id of Juno in GoN testnests?",
		encryptedPrice: "ZLAOLdxywOZB-U6_LlHhFnLh-1cwYT250CCSs-n6bqdkkN96sIFvCYWCjKE6uqKZ4YLam63lDbGIcU2sHZ62L2ezmwVHXlJJg19BECJRTIwsJFsvRTjIUILNj340uTcLFz47uBZepg6vacoKh4WxGtSVP0HBHgu-pskgHr88nw9q3tmjfg-UbeiFXR-kL13Xhn8SFryaEAo8yMiULwvhbWo=",
	},
	{
		question:       "What's the chain id of Stargaze in GoN testnests?",
		encryptedPrice: "V9-9WUC-Chn-RnfCPDQ0d5lRArAVpJ2aFOmCBusPOdalGEFDSf6luD0M0Tk-kgfwwrkUBC2HUR4AEoYCqsuiDY4-VPPwpIetH0r0YQpzf3Ga7xD1_1XOgpe64bnADMIuy3W1QGWkn64tnIuYY8C7kn15ujujlHKlN0Kn5r-Z9ZEL54IC-nAbM2rmaBsiNpB7Fd5_0qYRX_nRY8uJmepsMx4ODnHDqGjG5uGTUDLanew=",
	},
	{
		question:       "What's the chain id of Omniflix in GoN testnests?",
		encryptedPrice: "2aUoYR11fTaWToKZSnZbIQBXgc6esAfmoXba-NavG_90AGE96uP0VHC6lTseOQpzlyN9XqLp1LJ9b7dYx00_Wj7ZZZNDYs0RIBcrQkGcocbsotWp_dYKJ_nz8_9UzE1vVTH651d9RUzvLf_FT17HwltnTdjayTdQY8-TZAYqiBT88oRn7-Qg_XD-kxOrtSqfd_tEBwG_pnob01bLKQ==",
	},
	{
		question:       "How many ATOMs does the Game of NFTs phase 1 provide in prizes?",
		encryptedPrice: "xpM5PL66SGtCSFOIqm1BrV5LLiz1Fnh3cCEHuWdrPPhiW76kGVbhw3tdNsIr_j25OvdyajzQqM1Jea6eQ2f45wuTERX1JrV0rnuyVrEvKnTNxPksziUkIKCx03pqLpnbH33hVn6j6rfuq9cdB3FOIGf4ylQOHXKiy4RiBf-Mc_bmysgQ5-oz1a0s1QAtwT_0e4SbHl33EaHi-ILujyUqcubTx8YV",
	},
	{
		question:       "What is the proportion of the Game of NFTs phase 1 prize pool that will be used to reward the Contribution Awards winners?",
		encryptedPrice: "6jj81t0SWi-F50uTbWAGQ1_Vo4ehqfitdInv0xle82gJcGx0ltTHwIwC6djof_5JjHeWi5u0eqg-5iCQ7F8Mf2nhjiGb6MWP37s3uRqB3LiIIyHqn7-frQIz0_AEXA5-iY4DfpB2Tpy83j0wgjIkZhq7zCvJGaRRFtu0RT-bx1kvIxN1Nqd8hg-wn6LnrV0Ql5QvE6hX7jx17ta4Ib1jrQbRRIKkJNzW_rcqwsHkgQ==",
	},
	{
		question:       "What is the proportion of the Game of NFTs phase 1 prize pool that will be used to reward the Task Completion Awards winners?",
		encryptedPrice: "zyg59Y6O3LiSg7SEecqhQ5KmH-UT-2WHALSTAP2Gvni5vnHAa94NmkFS1AVgcWkcSz_6ltmI2jlYnm4tIzJqby2kUNYQdIK3C_jP-WhmONyE6kOMsrm_9-cc1SV8ynNUt1T9TSJP39b7b3876owAHzWpe2aZ8oaEkg-dc6iA_yH6u9QI8FcxPbv88UGDryq8zcBJRt1Iod10AvCTCcq4SOqCElE2",
	},
	{
		question:       "What is the native token of OmniFlix?",
		encryptedPrice: "rrstX6-jNJ5_JNNMnGBXSsCU9hMjr5N52ztbfKJ94BgVfL72nVntEOcrtavzAQ0uNgb4Jv57lAwlKU9TgNblotOyFUbRAonpXcwwCtdr1a4FbXjBhTp9lDkZFyP9nGZDYW-qeodRPY9dHlQ2L7wnTeP485U5vDjCX7hcVAwSDQtGb8kzRVQcHa4BlfQdfY4HzNBvYadFTNDDSy9_6-ciWSGjQ98cyD2kjA==",
	},
	{
		question:       "What is the native token of IRISnet?",
		encryptedPrice: "K5m2WcWtLzlpoANYseINERdzhqjayFJN1YfEzwT3Peqj-_1q8kwz5wthSxqegUyxbkHVQumipvSB9OpJGKOeuRFK6QF6ZhFEImAAe-oLfyBqk-SnFzTCqZe4u5Mz_kLjQMxnCNmj5nZp_kUX6wkDksBThh23Kccwgndm3DYo2QJ9RR1QKBnNp8fWiq_WbPx8jyKsGgRKYGUdFfSA9HCE14s=",
	},
	{
		question:       "What is the native token of Stargaze?",
		encryptedPrice: "sQy-pPbbdy00BSnN-lfbkG_P6LrkQC9f_983rBJE6tYPIdxR-n6ThRBRkLlydgV_sLl4PMyyJnxJPG5vDXqZVEoay05cPK6M_zwgGtq1kn9mPQ-oomllBpwbmmIsUGFECR4wpNI6KlGpGCaJKKz9wnYcrvOpuYu0xukk3WcQ6pVdHwFUQnzC_ZhQU8XYrfB602vpfOez8vy1F9mFyVjAcDZIoIk=",
	},
	{
		question:       "What is the native token of Juno?",
		encryptedPrice: "jGXAFSi4kt1dyFDkq_nxLQ8eXDT0RkNFUuRhmZeeNbAX0AYcG4xT3M37FZ2Uk0DdcPSzo2cYDSqLzQTyB4fdtI0_rxKBybH6_vefLhQ9hmIV415MK_aRJ9ALxVwdilxg-38BrNpel7Zh_aKITlZPvevYDONQRfqc6LNKV0-u7hWVOufZs7i7By-pm_nUQtBXoTsrHoTgkWL48JM-Ggj355jZoeiOieYTrpID",
	},
	{
		question:       "What's the abbreviation for the non-fungible token?",
		encryptedPrice: "co-jocwYgOKOIxBDUZbFZ2z-ssPY6ZRCWxqY3lKQ4B2C30of7mJXkx0LqMQpCDN2VLYDh2xYunmO_UbZf5TeAkL0prCXdC_0uktlcvzs0glsovFb6-02ZWnvGmE3NqIZP8CgK8XzqnRO6p4Rw4MGMz0kJRubhs4X2-9Evu80HS2YnTOvXmv_ASf-GyvKugTm_RkYNJXhs8qCCUm5FpilfggrQfn_7WE=",
	},
	{
		question:       "What's the operation to interchain transfer an nft on IRIS? (iris tx xxx)",
		encryptedPrice: "G8sa8QGwgpoB9ynkFgL1EZnQcwadzDsTxG8wHKjQkWGNCZdl9HUzRX1B8K_l0PBUB1czzEsawKgXbo3xsmPJs944riiNnOx2Sn-6ShbF3qPcmoKoBl66cbWqPOKWZ4JHEuwOxbYqU9N3gvDZkW8VhFJPehpsVNQk44xK-cLur0e1_nE6rtK5C_Q770JOpMEmENrOYsEFzcotF4M35cXZ4w0ykw==",
	},
	{
		question:       "What's the operation to create a class on IRIS? (iris tx nft xxx)",
		encryptedPrice: "G2eLL9PzYQvfWBBqHRyIeH-uJxIf03xiPvkPOfx5Gfv3kS5c6qD4WMIlXfD1pynLBDTfQIr-aRarMT1u2TJtbh6x6HYU8wcCxMCN7d8UNxodOYjbLMgQ1NfXiEwMRkFX8BBLscNYZgTQbXGsLyCq4AN6epj-LKhZSidBfBNTYHM1jvB2SIpJG-c-L7rYcygMRnvETY2kGM-4RsWLAE1iM7qkUCs=",
	},
	{
		question:       "What's the operation to generate a new nft on IRIS? (iris tx nft xxx)",
		encryptedPrice: "UYxIklbVcsGga29rZ8BvFG3c1SQ9KXJZEV8viMBeVDjDCxTV31NgFJ6m_zqh20IINaIxYBxBDksRvP4UfCOIU55GfMKD3NxsxSXT3MC-mB0SXRFqH2kmU-yfhK-S1j4_WEld73IazoNCWEYlDLkHo4coZ-JSt7aAojwB3MGxqB2Q991p6xpuTW7AFgtRxH6jPNjKXZtFSD4WUZql6v1_7l_gTYB3ggmbgk6TcNOS7qI=",
	},
	{
		question:       "What's the operation to destroy an nft on IRIS? (iris tx nft xxx)",
		encryptedPrice: "pFtGvVHWvDwHHNJ9EqPoUwzFqwiUhWhUv9sDSvtev6o9xvaArR6KgSom49eANMPuwBkysqah0I7Vk6IAbOfBUfDxq9cFrpmrVXAlRHonuCjTqYPVbQWm_IfyIk2CpQfYxFaT8WpATriI3bED-JDygBMp4DLwO4tXFWT41_4OOky_lIejhgNnFebx6mX-NyCQVwUjdFy9I7ilej6mgbdAXOabMiDFnqeL",
	},
	{
		question:       "What's the operation to give out the ownership of an NFT on IRIS? (iris tx nft xxx)",
		encryptedPrice: "TVql8wK4lHFqLNlKM6DZXAwVdJcYqhfEvZmjGDosOetEXdpeoRpFMwpNjkaH96Trfuk0kkbheF1zacnrhNGFTpf1TzrcnlHXlftVa5n38D314FLrRmXaPALk1yyekd1ixVbZ3ElyUxEgM1t2BGqowqzEcSiNCsAoPoU_mXHq4pak5Ca4xH0qvcS_K5SlOLVKhr4TRDAqbZpwWdBMLqpG1y6yZXC8qQ==",
	},
	{
		question:       "What's the operation to modify the NFT metadata in IRIS? (iris tx nft xxx)",
		encryptedPrice: "_14NXPnwEIGdAvD-jEDKFrpDN47CgDVvdi9_C7ck0zA1JO7B4-rt_GjpkpOAxD5ucGRJ9eb5wS2Di1crQG1d0X2eDppVcKc0jFcSaVhwpa-WGlamXMnT9JMT61nj-gfqP5-eXJEEpP5NoDcrPxFLuSDxii4EeGPUAUJIpylywYu5lGokaZimMEfzzTWQPjbkXmAOHWoNnTCzvgG_aP1J_EghlxlITHc=",
	},
	{
		question:       "What's the field used for an NFT to store off-chain metadata?",
		encryptedPrice: "SLskUnxY72pdvLQn3IxCFGdUCguHSFhp4WDjOepUVf_JIc6n6c_YB3Q-lejvgh_R4oAwNoAqc1s6YWxC5BJCOl8Rx_ke61jPfJ0LPG0M0QuAffeZEPBOl20EPjfYdLk9ABwwztR7AbO9Tc87csKUukhaG_23TO_vgC6nSm3mIrgR5oqDYv4cLNU4AS3ZjFgwjzN0omhmMlCwKlAXB6XBz2sHAQ==",
	},
	{
		question:       "What's the standard for NFT Rights Management? (erc-xxx)",
		encryptedPrice: "OV1yL09KI-nGJGLc5vUW7IpIp-6VDs1lRzeHMpooYy88JtbemSmbnbXq_50DtiCFLUlWdKDb8XQyzgjSVLGSkoHVuPLprXIXJ4ETBdBJZjr0DFo1WNzRzv4K_6yKR1lNIcOxniD40noy1apRMyk4YZ4DhvZhvL15fMe92JkL39nC6P6FseQLjEQywpFJe4tsiJQ-Stqyo6ZrBWJbRTxj0w==",
	},
	{
		question:       "What's the standard for NFT Hyperlink Extension? (erc-xxx)",
		encryptedPrice: "RG-Et9dX8uPi1W8QCGyTGU9BTsnevFWBwJy5eThEyJ0lR2X1MECbAm8Z6zRhx3o0IgrvywjO48nMzmz8zsycsmz9OYiM9sDlycqM8r2MUJsigBI4YpObi2eLPuozalBxrrWxvnduuuVU025H-gjBqSRY3oFJwypZaW3i4ijWALcZzcsbleC8igK7LRkPQKKVdXkTQMdg5kBVLSPf0JFvJEpxuvhpfKxDMQp0gw==",
	},
	{
		question:       "What's the standard for NFT? (erc-xxx)",
		encryptedPrice: "l1NaeB2GJd-MPGkebpMxndfCpIoxVNuRpdN0_foG-i1fhn5nhZwDm_RIJJNQPn1SatNWG4zpWCXs3qqJtE6YrzUlI8aq3orJ3UEIC5KCqWhYgc0LFhJaGl3xtlp9ym5JhM_yVA9AftY_gfbkFF5r6gCfSIHIjItFZezH6KrF8yUN_y1WwVZ3qEw2skRJsE_vD3t3O1fijQrRH6v_N-PkTn9QD9mAYQ==",
	},
	{
		question:       "What's the standard for Rentable NFT? (erc-xxx)",
		encryptedPrice: "Mv8Lll4f2Z3P3xnD6_tbE96MP8fbzAlhqlB43IX53KAHdGiPXisPkFjL8F_iVrHLNZDLXcVMxAQIdjmgqxyYqwhncItqZg52wGP3KMUG5d-l549Wc7pFpxTbsQvbWNyG19MaIWCEjix1jAgn6SKj95fRuxIUjCmTQCDuUUZ80dAj3e14-wzWIWv1DYnb2R66rhw2ACBSV3URB2QJZYM9nAEdkLL25MQLefs=",
	},
	{
		question:       "What's the standard for Royalty NFT? (erc-xxx)",
		encryptedPrice: "kwc9RlH8QySojwZYgYWFinjLAJZXxoPbig_yl0WdOQs8p_ptNDyPdOajFnb_bA9VDUoX0dDUgesZb0t0J7IlFqPLSjd4teA-Qz9h7KCR-Ot1vjfrrDEP0NC50FkgDiEhtkIi_QUGNOvKl_YJWAKK_mI2Ve0bp-K5zkaSY_V27xKPoAq8FIyiqPoxQaoxIQ5-zYmG7Uc3NYfAp5GfzxMqShcjFwIcvHQwh50nY-PNu8LwY0gfnw==",
	},
	{
		question:       "What's the standard for NFT Licensing Agreements ? (erc-xxx)",
		encryptedPrice: "iHH9zejw-EP_HpoLKwvW5cB9-_CkX2kdcMFeMlB2qhBa25_viw7-UG1eZ6jEJUyG2MYh1rr1KYMeovf41t-1234jmuYCJTnIqkUJqXjIf_EcqNw23HnshjGuC0WSNkCqv9nLcjubXiT3ju2VtCJvOuc8H8ACiVwV0IIljnj1Un8UTnBVPAMLwBuj6Dya894Lny7Qfha-UHamsJJNhI7iPY1GsOhu",
	},
	{
		question:       "What's the standard for SFT? (erc-xxx)",
		encryptedPrice: "kdab2utMBNEX4sA2aK_3X2aPwp0S2800ejy3Ag7m1dNEfUfOWle8ny8vcTv9L2ZBGK3YbOdxsiNivwR8fUeuQdokv1FObehLhttlubE1y-fWFpgoFVZMGQFcebM7DnC-Tpxn6aD2lp5RzEtkNyJSCI-gtBTX0Yiu6cWN_22NLc-pX0TnvItgAFCFRBLIDYogVE2FRF1xRvJlyPHvufrozFugOy7uGg==",
	},
	{
		question:       "What's the standard for NFT Consumable Extension? (erc-xxx)",
		encryptedPrice: "76swDqtZFoIWfak8d5P00t5CezVl1IkLMNt9b4TCP5Lc4QDNQLlr1hXN7L7Rmh2nOpVrzeyYwRg_uCAHqUg-agVA5ygIIrxsxFZ-WtZqUeiCQizzpd0a3aAsdrr3XnYKvR5R62U4EcNWyqNi5amR-vEXoyU_RA8n90fNlPELG6564A9D7NliGqZDD83FSkIk2g6une6cW3-dqSK54Iry2VHKroxEjYdzqw==",
	},
	{
		question:       "What's the standard for Hierarchical NFTs? (erc-xxx)",
		encryptedPrice: "lELcFzAlJHzOU7WCbmgMqsr0qg3ldUcZtOMr8zOtt6JcPhtpGbB4hrEiNGlwJyJ9ef3uXAqcKjN4D0VG52ff2e0rveNiK-4p9uQ8uD_2nUBfqRkGwqgbu0n_hk7YmEGmY2k2pewggfRJoLf1y2tUzRu-SVRj9JW3522Lm_AjK3lyR0iHYkwpn7jHYqtBkOTRqEnAi0wtKE6O8zpuLSWA-jqgAc7j",
	},
	{
		question:       "What's the standard for Shareable NFTs? (erc-xxx)",
		encryptedPrice: "jnNq3n3t-FWLS3fhUDsKYbooB-Dc6htUldfEommEeSL9MJYyKJWUZoF_gI1x5JuW6y_n13LhePfcEPxR2DPADGzL8WliD07JDcAeXA-vJS49e1mO6tcERjPgJDJ53oeqZiIm_lXmT5L-1zo5DI6N6X8xfMZPVxFM7Wlm9Jo-_2bP74NxeUyOtGin0m8OwZfNMtCO-LpW6_sNI7auypizcXN-tnwp5x2K",
	},
	{
		question:       "What is the repository name of the Go implementation of IBC relayer?",
		encryptedPrice: "kQaEeKpQyniVQLuaWPfzonVW7DrT36_zFO2sG8Dc532_SNe43R5MuL1EKnCuwG5U5_bL4h8M9mz1KKaFgUwlv4QZ0yYBTPMnhHLO6h8ACa4jYDDh-uA5cPbxY_B7i1lxmTJHsE0rtucFn1Bod_qP3Fjonst-zE9Dzj2skjtbN57Blei_TpDWzn-BumIFwOAaJAAYklBsVTru4RghOTWu_n1v",
	},
	{
		question:       "What is the repository name of the Rust implementation of IBC relayer?",
		encryptedPrice: "S6iut5F_S9frGH1kRf6oAL4VmqZ1U3FoEBFeRKIZHEi0ElHzj146_CTV7ZoughSA19__RHGIcF-0O0RBxMatrUo8ncMmqdzTI48-P2G6fkBSZDkeMzA0hWpbu-7LZoSxd4DPFWlFCkEUX6gW50tWtItIYWcoz-fsukCRU0HCf2qiI85JdL8fHgWrnAF5TQfg4bdVDORedDxlnpWvkBM1IEsLC0wNDxf76R0F4ydjP_npeID6wto=",
	},
	{
		question:       "What is the repository name of the SDK implementation of ICS-721 by Bianjie?",
		encryptedPrice: "HRKPEwmY3IG_DpZYq5BWtvONM9D7MPmSojwQ76GUGkZCXcX4p9HN9O3uuf1odzMqeaG3DXEANN4jf54EKOcmwtW9mwxZgJ5mQ2qBlDTrBPssVA02GWrzAmx0MrBY7lwlASdekLlIq82iF3rJBjkxO30-Y_2tem0op7HoGu8zVVF9wjMez_rP0uvSTfHiJerzYRn8yRWN9kE_q-C0hYuqzc6Fa7yX",
	},
	{
		question:       "What is the repository name of the Wasm implementation of ICS-721 by Stargaze?",
		encryptedPrice: "6A850VIydbQvi5TaPq9h0ERz4TzVZZvpdOPboL-l8uO9CIe51RDw6gACzmDQJlGQIJRmD-hnLS-f35kvJ4Mq4q__21Xq5bLZP0QOg_T0ZOLwJYfizFVWWW6XKpv99dZsS3TqeOARbFP9t97awYAC79wf69B-kdCmbm76luiNTQr7eL-KwkavSAdSAOIMdAKg4BRW1TzLv_5LSSaiI8rbZPLCX7bd",
	},
	{
		question:       "How many general tasks are in GoN event?",
		encryptedPrice: "Q4z6LB2Dni2xZxdaCSlv36AQPbSkQK-HHadS3Z32zj7lfdTcdqlYhAVMI4LNlcQ6GIijgxkLIf3-p2w7Lvv0rTVOk-tseolEbPrpGfDbI8t_xY-qKenTVmBHK02JbK2MrPKpLT4CjkO-_2QPjaUQf7-3GQdT6c5S8IXlrdJ_fdPsLAg3xqNuHyLTU9JDSguuED9TkdkpqYRoEhZfi18HJZ8COrtWRjnfw6bfqk4s86w=",
	},
	{
		question:       "How many game tasks are in GoN event?",
		encryptedPrice: "v201JWiL_yZo99DdnxI3qNNxtisjaiJ62hWtHqmOcaKmp1itT_ABIe69Fc8bJSrWGbGoBCi0mf-GjJx8Tu8KLWFKgMWWdNqt5cfzKRBh_moKuYxxYtIThQyqt9PWz-G3r34Z1t_rlfpWim8wo7UtWmCmK_zL_ewho1BSscVkMTyFVkR1_hRnSp2IFo-eR3BvunGg_uU28Y6TsMMW",
	},
	{
		question:       "How many challenges are in GoN event?",
		encryptedPrice: "kCr0i07tgb1S1GMc2ZzzAtk2-HbKFjALIVXHYXAUkD8xzVbsOHRMk1cmxmL4Al4hGGyiUm5zDgxVdRB6CVFyOlH7eU97UivaQ_FzhZh57KhKuE-x1WkCz3eU90adDpVjFNBcMHbtxnmvI4AGZHuyDyXDKojldSnXHwCay8FDi2_2K3SUmkA7bN_dsHxdGayaVZWKvFIpoRK-nv-JX1LvtovHn9bpNnEu6w==",
	},
	{
		question:       "How many awards are in GoN event?",
		encryptedPrice: "tlgIM1cC3Z5_52qGGMXfRVkLieVk3zsxIICUT91PRW8Z0ZavjIrU6-7Yt-_IyN-mVWkbeJAZVWOHDqCEqwc9xY7izuCMtx20iK6--wcJHf4OEjl-peUIBfENf4K0-J15oVByZcwOLO5ik3oQV7xI_XFnDjfqWyLRW3CVSklWbuBl3FNUqnuFknvHUiKyaIJeo7kIDDh-0k0989PyXHefY0k5UuxFSA8OMRoM",
	},
	{
		question:       "How many round of airdrops are in GoN event?",
		encryptedPrice: "Bdq-YyYJdEfG7FTyK2pmoP5BcejNQop-OrxXXlI2L1LPtupJlBlu28xDLw9ywlrhkSWHH0xtARFs_5MIaT-hUOhDEDWIMXeJm4bo6sDMxyW06XvjfVzroksg8dF_8sy47EC2G2Ke1Aw0J8OnB-zQ7z0emuUEaJGOB3i_ovxC6QjSmntoDQNGwEGSYah_4rE0o61Dom8pora7-lvgLk6dax0bKiT17A==",
	},
	{
		question:       "Which Interchain Standard has specified the interchain account implementation?",
		encryptedPrice: "hdyXm__AtREK6uJg-l3imjAOmdopoOjPJHi3NlnyJlLgrRIC26DgnCanfib5SvpHJV5RbnDRDoSPCdCciGNTUjmUlB2BJc3_HspRbZ0BYP1baYfPZlz9TWknubHw0wwqqp-zfq25jeVSsKDhlAt6DRk1eNopA46r0OvqqC4ulhLYwiCq1X0FqdI8mDUkcpvdPdHcxHUWf15RdXTKN1S4piZ7HQbvMsyajQ==",
	},
	{
		question:       "Which Interchain Standard has specified the atom swap implementation?",
		encryptedPrice: "tc6fCxdnWFDV0EsQlyyp7uZimY-zSvyBnGL52et-YjR-mu3pCH06MNbeXj2U1faHIEz61j5zal5aKASnxehSSi_L3bADekOCJDDZjr3s8CEuFtwomUgmyxrRD3ZDnLg389xlt1C1qthMfFgbj7GHcJFhLuWPVpCFNrT3rW4ok5-B5tZBTgsRWznWX0iPtazLDF4sRYQYy3H0qGrY4DK3J43YWsmWFQ==",
	},
	{
		question:       "Which part of the core IBC protocol serves as the conduit to transfer packets between a module/application on one chain and a module on another chain?",
		encryptedPrice: "mxksB_GdYEEXjRPm1dpnKtRmuaM8ldufkFyrQh7ywZrzH_neZxH-RIBweAtkXBpAP9ZKwKkzWHSecxWFU27YUpeqtDIubBQRfnfOrV0o9dIxoXQGzjvYpZ9NdGADNtEwzN_vQ3JPJtjUW9McqCNg1H32d1Rerr_oPJWr55E8n-nyLMnAjsTV-3FR_FqjuR8SlzUorNqTh3gv0eaNd5-snhG9GYn6PEPFrnMOqp3T",
	},
	{
		question:       "What is the particular kind of identifier that is used in permission channel opening in IBC?",
		encryptedPrice: "PPbAfOch9CQRDrVBpAwd3u_dTs2ZHlMqiAmTy-wzk8cb3LhnG9y9Jp3om9DXJuyOt1EbnpxcN9THt_h3UFbTQY5c-DNPFXBS-q3-M27S0To_VXuexSEdRUoDDYCFNDpsVxo-f1P8Mr9PSCC5MP5YHt3__PnqYwmzFzoUzVkKBizLNa1Vgjqr1uQg75vpe-7qWINZYFtVq9oBmIIG",
	},
	{
		question:       "What is the name of the protocol utilized by Cosmos for resolving blockchain interoperability?",
		encryptedPrice: "Z4LRCBFJxhGMn_tWyUqE8lrHvnwwyuyOvJYCSgitmhH1U4_bK3Ux4wyticL3IZobpnVP94tmsmxmVcqWjviGXpNQZVOtnrN1AkfEWSq21trrW0fBJVB3g5iXG_Pwwx6JmZTyyEwFdkK3F0ueRR7xmHNo6oPRDKpqlXTUQ6b2Qp80NbstUXyWAZG0V4CwIJYXeUKokrq4F3fRRFSNx79V1oAIlTxc",
	},
	{
		question:       "What does IBC use to track the consensus states of other blockchains?",
		encryptedPrice: "Izl5F2nU8LEmQmfyYUVaOG_RLWWhBVwUE_Q4E1CmAs6bzmwLDVuFmbOUHtk2m_5sMz_AErcxg8co7Q1gLIXDDMZ5iqRo26hM9DUU70jD8MmYn1nh0rM__jOYoxvyVtVKbpnGYycH211UBdgju2X5YOqix272OdGZuD7_l9kg1VBadFTFyspAIa1jj5CiVHGlu1iYsRnzaWzjOZ_8ZjFcOSLoB-36s4ymnx5vf-zO3g==",
	},
	{
		question:       "Which Interchain Standard has specified the interchain NFT transfer implementation?",
		encryptedPrice: "GRq1u2jXxtpUaPJLRNBf7FAo0sVm8TAqw7feQ4Qi5DqGyNA4Nb9WHkV_fq-Zw2M_bstyoQEXpUxvwFrIPCKPlJ2P1njJCbAhV02lhxuz-GybnIcHosbQ5XPPRUG0dqJxcu1r2Qx1juT8KMnkyk9qPJguJk2q0j3oAbN1LNNvPV0V9sO22lVjz2D7RM-SkaVCrM553VnI1xHesbMFZec9jdY=",
	},
	{
		question:       "Which Interchain Standard has specified the interchain fungible token transfer implementation?",
		encryptedPrice: "z8wnye_NZEtSUxXy-DVazATsPbzqLfCg7rIN5j_N61vMaFqTJuWXZkBFQ8tjWDtxxlGVICyMPGot5NOf9Mc6ERAAODme7SIfb-qo-gI4CO4ogku5vA0pmE4jnDEXOP3FWqDgG491ge94BC-4uhK5Yx9vyfeOiwfv7EQA1vtTZNNZLiSA7UZz56CGT_xcXaQ4aVFI5TPn5N-n6z6GE9pRv3X9x-Vw",
	},
	{
		question:       "Which Interchain Standard has specified the fee middleware implementation for incentivizing Cosmos relayers?",
		encryptedPrice: "2Coe_9p4vNNjomB4VFMHwJwzj9k9PbflckK9JP7XgohOgfW6YXDFN6T5TC8ysHymLs-RE6b-CLCkRcSXkwpidXDcwkallfxhJb79SetwK38BGW82E3M_1tzvOwA4I3rgFatXMNuxAAi_UxXeG5JwaEP8GVVv4cFzsZc5cuXimbKFeyOru9mQQBlxvMcLavsA6BieCS3EWg-GymiafL2oc-fMTa2fEW_jmx9f0WcljDdDSg==",
	},
	{
		question:       "Which Interchain Standard has specified the crosschain validation implementation?",
		encryptedPrice: "iWCpmD6NXjehat0uBWUmBbTiqCkEtxt-fA6xURAkx1_fim2vmDo9sdf4xsryoC4syNbH_JMUAdyo_Tzq7iyl-libbQ1WJKcJ9R-brv6SXIoYwOcfAhtW5jMDs5E6c0GC65U0S58RVretLnHd-8--Ao9Q_yo59bLu_fy_ZPrFrDbxzCQXUHI5vpiBMXAGve65xLnhRU6yb6oHn3Dpw9rL-79oPzea6rA=",
	},
	{
		question:       "Which Interchain Standard has specified the crosschain query implementation?",
		encryptedPrice: "vMVDhRqPiipJYSNsQUGQgTAP31rU8ydOLEnGzpjRejraE_OjoYcx9R3dW-PJlb997RlahJeYTNHls_AVNbGH1pJXg2vVWlmDNnI6JSN8Nj0i7Ew0PYijZqgJdgKjbYthZ99WhP8RJW1NRKtG2_0AYJbWPxfbM92afbBE63xHkfNsGaGw6aJ739E1_cRQ6qkMZQcNjgBMw3CUs9krq88yaCPpCuEceIpg6qHEVYJj8DkW",
	},
	{
		question:       "What's the default light client type for Cosmos Chain?",
		encryptedPrice: "c_0U3vj6TKjB3blmnvb6IHq9a2nODCyH39kFAxAKOEz_FE4ivbihjBi_Oz1xB9kAuBHtJ-iZSJdS73hdmL29qM_6zQ1H8OKboQwQAFArDly_u2GFeBT-0z5BEDkgN_S4s3f3C6q3SYexgYNzkkygIsgIBYuSCs3cbSi9dqGDSvODi1mNYiuHTE1nNfWU_-_PEJjaQ4fkHYb8ZjqP8-IQaKvgfd_LoQ==",
	},
	{
		question:       "What's the recommended data format for tokenData and classData in NonFungibleTokenPacketData?",
		encryptedPrice: "VnIIcmPTCYBsaq3HtgbFMonindc38e54iI-w64Bm7fKGbLTAmBCoAD_x3LLRE3WkhgjBRFezEWn7lMr1KfN-WntzgJEyIn0lVWrY_D3MUR1UNvQfNl4iLnY7aAV7BWvvkrU0QsRDdvYNZe4yjsZXG1qslgc3gmhUz9vYJHgYA3SynfIgZPM-mfFaNCemBG73gZptB5Lp6Q==",
	},
	{
		question:       "What's the required encoding type for tokenData and classData in NonFungibleTokenPacketData?",
		encryptedPrice: "gct21fiZMGCnZR1HJcxR59AF4uvPpYvOrMIIX5e6tjNoKCs2P3IxBVKSoaRkzaBe8NjORERLdP5WO0Yt_dW88xlDHJQMUG3aPSQosLXzLEdXZPmrWoWgg-qUhOYSYyItR8I9tQHx_R5OFKM3UOhLEzsxIssRS08N_RvmZ3zig6_lkbawswCAKFceph77W0pLvOTNpINM2WIl4KR9ZMaKRIHIyOnTPzhTOvjBzTrs6w==",
	},
	{
		question:       "In which programming language did IRISnet implement ICS-721?",
		encryptedPrice: "JDnMIYtB8ONpZ9YcdaYydQ-yVY93Ea-14HTvSLy9HZNJAPc3AHvVu4MC0zWEcwqW-Cg6fhb5-YgarGc1X2reJSLFnlj7BKQmTeOIOGxDnIId4-WfSMdeOfTVca2nFYzRqLECs1wGUCiMQpxs4iZEW6WuMkrfgfiKYF74ugQukW_lCiLjRncP1rJjW__WkzrMaUoPMvoTPtMWcVLz8CzfnW2h2-C737BIxwnCnDunARty",
	},
	{
		question:       "Through which framework did Juno and Stargaze implement ICS-721?",
		encryptedPrice: "-qelSIdHkxJtZvlg4VYEA4TH8KfT9zebZPHVCJ0OKxfTAhTzcbtWywgAxNxhuggYahimhLFEQtiiylOhVF4ykIczcHyk6voRJHP8syOXeBbuEO1uaF988NUT7_ID1zVa1GAZExfinzQqvlilv7hKcMGbrD2hKDyqBZGtwv_3pi9A7xDURKtz0Qftp77qL3RkIoa1gtElSQpa1DiVRCueZmhNdbM3",
	},
	{
		question:       "What's the repo's name of the Cosmos Hub?",
		encryptedPrice: "m1fw5NUzWMqk6QSYVyJibtf5Y-La1R-_8k5Cvtb94w61qGPbFzaoTrT6m-CD5XYDN29U3y-TLIDoFZf_a4YxZ7hz14KpJXeeDiNlSjWcuDBg4N_3kP46vSEUS9Eu43IUMWFQOc1NyL9vC5YQu-BHFekxi5NnIc8AQdCNGNUS365C5jRUBLJ5svBlwpoNs-ggxijnFoP1htRK87b3mih2aWhe3Dx5K7AfKfc=",
	},
	{
		question:       "What's the token symbol of Cosmos Hub?",
		encryptedPrice: "6U-WygKeWF9cZ5bRCE_JK8ND_plO4t-ICZo9fgByRjF7JM3bG7Sp49cdB0XIuB7LFkdfu3HXqdSJeUoVv3ywtyXkvVi2-o5RVW9DO3uS-by4TciY6YM8UaJOaG_0KROx5FdUzjK9PMR6yjtTTl3mn7R36LWzVpzHzI57xOT4xFjkCNF611UbQrnmmZqbQtLu85Z5MGLYaldNIzatNx5Ea6XLdDd9I_aMQO9CpAkJui1wrQ==",
	},
	{
		question:       "What interface does the application layer use to communicate with the consensus layer?",
		encryptedPrice: "Bp6OS4iamWRGmC4tzf7_pgoBGM_2M-MXZfD8ipKzl7rQ0P2O8qXkEURotq4FTf5jkvLhbLLGkOzlBRjURU-OkIsQnskV3uaEF_t8th9YYk8MHGXO8NoAA9ardVd_UmpqtwsnedJTTDv9A2lmNhhzj8h5BUm3Ct2tsPmphAI7ImHVNeR1bANu0EK2JJmAmRJ3pFYNbL9E9i1lX2sjPaBDejT-4tk=",
	},
	{
		question:       "What terminology cosmos uses for tracking resource consumption?",
		encryptedPrice: "Y5jx_W4TxSvYKs7HY8WXERJhNeMf1SdF6rj713TYVmJwxtOS8SiM2uvmo2Jj_8ZVYEonbiLEN4YlYkDvQ8EUydTws_vCnEDARYOv9Z86gRP98G2-W_DGIMbwBXHhOualFyu22bT75utDCmeRweakkIWqreWyzRSPgDoY68vtZoAk0qvwBZF69qEKkxKHgascc_79oJ1Kuq7DkNI4fzEXJ_83kuzPxqbZzQ==",
	},
	{
		question:       "Which role is responsible for commiting blocks?",
		encryptedPrice: "GiA6k4QUEn1JFoo73Z7XL9ZiKOGiWAGh8sAvK2W9SJzxnWPmLX9WV4l9tizKxKsPixGfvdZKRLuuNXSx3rz5hT3jueEqXEt9nXHGDYcklRKB3Jvcynwp1wWIUqrrqAwigiNvVonnxfTIQG-OXM_1UveXrIM2no22EH-bJDA7yPHC-_rfS-mBSfGMFEkKwpfHkgYAo02yvugZf9UosrKQ4rAh",
	},
	{
		question:       "What's the process of locking up a digital asset to provide economic security for a public blockchain?",
		encryptedPrice: "aeDAkSZUwy-Mca0HNLsshDdej9e1kfECD4vyw-vs44hlhyPMiAEbiRGQzQEly1z2tjdGDHC67KhWGa3jcVqYq3FdVc7BZ8w29uzXjztRJl5B49I7iV8LSt_OFW60YthiFv2Z0SF2YTyj8HOzHWfwR744_KE3I1_x2Wl2qZ7cwX2tpK7SlpHw6Cuze5Hb-CCmenRJ3W7ybfOT7uRc3LAZzH6wbCswM6srnQ==",
	},
	{
		question:       "What object does a user create to trigger state changes in application?",
		encryptedPrice: "JVzs5-7l8aaP_QCQ405ijqyIh12DchNUZ3udG3LiV52V8NhqzHy7d8SdMpjuT7g_zjIJu61DAAz3bOTq-67xb6u4QQoOY5L1hwi2vMgT6aP8ILcTp9n_SOjiiih9xXXnQPhXTfY0SqWRrQm2TpWbte96I_plgyK1caGMx-aBKUH2vpJFQaaumky6yUS28GH-sqaHJni3v3WMJDlmrxdosk2MIiEBYA==",
	},
	{
		question:       "What's the one step before including message in a transaction?",
		encryptedPrice: "ApbL3_SPc2EmCpn2CroDUIqk1Odj44i2BfE80uDr_GzmHRgjqI2Kxuux0XCIvQi5lLOIITOSj9hGJfCy25FbjN2nxRyh_GcGqw-vc9wMEV4iC5BLaiOd9aVG7PxYuFCkko0XdvwRtALg__0hXqEVC7AbsDe65R_F04Iolv1Gvxc37T_OwJH6RJAA59Uw9E32K5jtaxyC4-j_RGt_GwIh8VFvrg==",
	},
	{
		question:       "What's the terminology for transactions being sent to other nodes?",
		encryptedPrice: "8Wm-zqrYZ04G4OvRmfeCM4oazkuKs8wAwk6TtYIp1y_2nnfv3KAZGmz1ijCA_d2l8R8dnBETK0mwTodGULAnL7y_lMNoQo67DfITy46qpj2mdnS0u_itMfS3JFCSlGO_LIjnef_zZQNehid8xTcBKHvpvxUXPG19b9Wem-CCflqmn03MhHufjmpjXBAeEiJRuOReNAmkOZ5StvjVj0c=",
	},
	{
		question:       "Through what mechanism can atom holders submit a proposal?",
		encryptedPrice: "Y8U26BvCarytYLD2WlhO1QSPLYjL_stOC4XwKOdOhdk0V063dQyWENEWtKsOFb50PBgzF6atAR8MMXZYm_iOfbFH1HaVz8K8-x9mElU-BWDdHtVYBXwEbZYamjHQW5cyhrxlmolRK-XrYQoDr3Vim3wtM_LetqlQNRBo-C1zHudtIhLOF7D11vcKxZ43Vbc_JmVuFp9DvWOewPZCWAD4",
	},
	{
		question:       "What's the key prerequisite for a proposer to submit a proposal?",
		encryptedPrice: "E0arQF9mU66dnjHPVY5es6I761rNqvdcvC-GwtabPY6LkWsQ_NnpVssjhPOk2Np2Bide-tI3oJvcAoeubKI2YbUTq1yDxxFcUnObk6GLLCmxXOW9MnP3sdYs3SVPBk7Msk9wro3Lf2bm0OJ3N4DruyHJZK9oNl1DW8Hbq_H1FQ5nyC75rsZFoGWKqDDRCT3BTdLLPY1bPHr_xEAMI2U_eT8QoFapa5dkIQ==",
	},
	{
		question:       "What's the operation to impact the result of a proposal?",
		encryptedPrice: "cRn-fYma-b5a5ghL0R3z9-tUTrBR4-9DDFC4An2CkgJT1-5fmRfledIcxtD5mq2MoHVjh5uKYSK1Bmmd1nK9C-4I1UW9Ezu9fFrd68Qlfs66czvdXLgMQQgsvjPEWpOTcyfNPDvLwPTpNi99JwPG3yx0-iuxnUeh75Tf9KTBbZmF17ebRgYie7kDNiY5ibb-IjNgpDa1nOeGJkzVK2mtPKrFKyXRN4Po4jyOJJLJxDKRn50=",
	},
	{
		question:       "What's the terminology for the minimum percentage of voting power that needs to be cast on a proposal for the result to be valid?",
		encryptedPrice: "RglUuAmiV6uWwm639MmFMZqlWki_UHuO84pqwkgqwxRLUbBTf0G9kQhIdo8VZDATatTIdEIVgRS52IIlv1eKU8UxFHeE1dok9ov9Y-hGf1j6y877CBaAtJU9YXsCm-ketBYKU1ZQHTBr0F-LpEd6tjY82YSOLbxOZ7Uj90K2o-MC_MsMAn8wFWOCUvK254HTBuF7eVpRI0Sr22_bQmQqh0P1gGw=",
	},
	{
		question:       "Through what tech does the Cosmos gain horizontal scalability?",
		encryptedPrice: "NJ1Rpk6p45ko4jRLeHg5ta5xcDWj0Zdat2erlsyKP8WOG0j9oIBBGsWFLajvALMMqgEsTThRkUQRangcl4dvrbf2USioqmlzoAMsyp4oqyhEMtSv3csGfrAoUvrBJiHKZdx4yywm_fUZ0RO4imUo_jHbE7PGs3V5yz972gpToR1M9chTYXAiQz-eBvthkhD23_0W8ylT13h4ygJeP5G3HNk=",
	},
	{
		question:       "Which chain in Cosmos ecosystem first implemented EVM?",
		encryptedPrice: "I1QcQ0Wd4wfu73nl_GDZJ99_n6W6GL_b7_FkhHnK96JhI3WiE8XKNEG1HDnmIcBCtP9Yt6114qUxGZTupI1-X0jNnWAZFimf03yuk0CrmLhmZpEhMwiSb4PxidPchq3QTqMEziFmyXl_qMOe5qf1jomTk6EWmnSXV7fs2s2fh-2V5i5hLm7QG5phAQhO9R9Ulm3iuwqNvo-FD3TKuVw79oOG1mTb5xSiZwQCUvbhFjpOvQ==",
	},
	{
		question:       "What year did Cosmos Hub first go live?",
		encryptedPrice: "FBP_1TEpEzBRqq2CV3Iow1JV17LOHLyHUwc1ViGFMLHJsL0GX_PYE0EYQhdnq8UfGWqkQ43nFEndvqPHNta7S0eE848uYldCEvB-qHNH4KcPzZNgNtJbLDkvR_r-L2A6W4D1oHjKX0hIf1LvWOgohV-7XvW1HCDQtNXE8fvVg7WBlQa7ywe6J5gN82Q7GF8-xium1MwefFE4RnsAbbKbvr3N3DM=",
	},
	{
		question:       "What year did IRIS Hub first go live?",
		encryptedPrice: "xiYzhYfyN1ouQp2IooDyfWW2MwrI1dnyVqBz83y3A0Ev7OivsppIokpZt4TSpREdl5pKdIk8pOp4WULRUXxW_D-F6sXYHVRer0XIKooxskyADzQetGZREWMwDkqmP2s7Xvv7Eokz8OgfZYIoVHJhzPvxOIPOfBf8PAcaz5UkPijPm_4HL3uMFQ-bsZYDx7w84NiZ0XLsvkMBKUw=",
	},
	{
		question:       "What year did Stargaze first go live?",
		encryptedPrice: "mTivzdA4GJ7BthTm0vzG9_lP0QZpvIdup0bMbfXZk3yIpx1E8xxZzM9s1OyKrHZLLt-CqTnltbE_3rq4PuJ9eT7y_coVA_nblb-844TdP7sMz1XvG4AnXTrItMesaaVMTTrKIaOQu5pIUe2265PCjkaJb6w2s92h-MM4mD3xbKg2F4GWjY-GAHOV8LJU6y-CC88oCUcOTdFytYmM_ItyxG9OfGUdkDKNyT0Sdith6_rpRnU=",
	},
	{
		question:       "What year did Juno first go live?",
		encryptedPrice: "zaKv0Me6nFumk0QwUDZJ_YuMVhgTbs4R6U_oR4QhvovgzHwN-HpgMRK1FSeCFjd-qiYME2_uoIDkBlWAXm6Z6hRZfB02A8D34EfZXXggDR8CxxIH_cSFvEiv_bskSKXvSDNLXlNve_pDSXYHiR69if2Wzz-tId_TbnfC1hL_5w-_3rgeLoVG1wuj-1b_lbIP9UL7txnZBEFNyJJ6AcEJge7mXqedyjk5bmA=",
	},
	{
		question:       "What year did IBC first launch?",
		encryptedPrice: "eBUYyoKq_dF8mdHZpXFt6HwI9egmjZupyw6M7Fj4ZtumrVGg8zf7u7UOseDfS9cRjnEihAHyD2YjDNbxjeZYazZzsFfJaXhH_DNO9X2JjC9JCwk8fhghcflmCnYPPSYZ-ZBk_3mQiiikV55_z5HRcoDA6GRQc3KKOXeiy8dJQyR7DmEK1CUYZ0MfIUA1wNt9qGOGllCMaXLxCgjFV0xHd3Ed8LiK5w==",
	},
	{
		question:       "Who are the main contributors to the ICS-721 proposal?",
		encryptedPrice: "PJ8VcSdfCmAdToP2pmp8freeFElwm1HyGW4IbG9hlLHe4nCxcLB5mvA_7QPiSdcU5SjQ-sdHTv0Wxv4aD92DGysffdPhobuM5KcUOX67BDJAujRuEthnD_K8wr6zZikshRB6ukCqeYJhwlr4dv0NYYG3eCzr79W-a4qYwbiOFrvr_NqUvbYWioZgL76AJnwDu9zSu-JgxIynquUKpJ7r-FvivgijbGw=",
	},
	{
		question:       "What is the group of chains that have enabled IBC called?",
		encryptedPrice: "HN13EYuG0HrBFha7BvAR3wpiSoOqb6ARSMzB7_akq6MIkahFGECmPzj3crGoRg9OQA7mWcogg7QWD8BpUqTZmoIHJrabc2uOIuHkr534Q17IpJ6To0_GVl5fgrgR-9tKM4g17WPNpWVMTxfwWfHBPsE3wp5Sj0EDTQedVlrnfG9CIovHCGJ1XB5qNR7e4gnHqjS_Wv4K57LVm-4eX9NYtDRLbSUO10I=",
	},
	{
		question:       "How many chains are interconnected through IBC, according to mapofzones.com?",
		encryptedPrice: "WT94V-CsFb5pzFFfq5z7VBPcH-7hxxWAlijbgtIpbgEbaWhDz_ZHQnzC572PuM6HololXNHu0cvN9EJIH2TN816y66Mxgv0qGAwx-CQHCExxoVAhMzuNmoDrckzhMpR_AzOUMBm7IzDYH8NPZx-0txPI09RyJk8dfzusMfBXyjXawk4pZV86J3-_oErQekcTJL5cRSKcGq1YRjq9yXF99us9",
	},
	{
		question:       "In which year was the Cosmos Hub launched?",
		encryptedPrice: "z6ruhNkI0gLlAxm3xkTvxg1tHg3jTsw8tMtpey67wlQa_h4mJAEsqkFGDN9qCRjM52GgTxL3jR3wdt_Tyq6wF4dMW8DdJSv1GspTDjFFxDAJ1QmemxQrLGy49RDOfEBVdl55SAEGbI0AvxgDYjr_WrZ1fd6pHGMenVx7l6uDoBbt9PpeHF35mwaRaVSd1OgmS6IY16A5A6w-YpOgLOBxEcpQs3y1VyvSGC8q",
	},
	{
		question:       "In which year was the IRISnet mainnet launched?",
		encryptedPrice: "Hto6l6I4b0rnnLPuIFz4ErXK6C6RVZFeiIv-B5tCgPUElPeZW0NZUNKhVo0_oht8fIzns-E1lyDzmXy7QNPRtenzCrWU7Jzjq2SdXSrs8k4KpXKL7XwEA3dq6y0Z0-OM2pII__0_ZRIkiS-jz-aKw6GmoAC5o62lDz6ncejUHOI2Hq4tAVJ_riU4EsUsaEgJJRBAJFyfaj0GtijPr9ZOYClVBd4=",
	},
	{
		question:       "In which year was the OmniFlix mainnet launched?",
		encryptedPrice: "X9o8gbVHV0qYt0dU7sjykt-Dm2DzqcS1xeY7mmh9TBTJM0tMAk0KuVBb3xmKb6yUJvnzf8TfSqj-3jBwpWsxF3wxsUk3OepNWn4WgMwwz6DYfsbIHvvXV7wJPLhI0nsfw6S7xIQHUahcPUXbKRJ8FYuw3n1zCFvXSsMtvqglbuWHaBRf_OT5EnK2eCbsn0YEZyXtUVIlwRce_LQ_K7U=",
	},
	{
		question:       "In which year was the Stargaze mainnet launched?",
		encryptedPrice: "DBvoSqZ0v0bJqIRPDcAdDpD8dkiELgv3jiyblsRc_9uOnbFqOqYlxjAy3VN-_mY1BoorpJsU4mwLbLUGNagynmZJCRjPUUhpZcBvrNtNUpdlztnpoX1C-kV6iG5tuH9h_YiSq6d6V8UoERRXVFcOisR1jY-f_i8MtZMqtNoa7STCSl_W3NznVBMxaJUBx9N94a4Qq6HPbQR56n2HG549Dg==",
	},
	{
		question:       "In which year was the Uptick-Chain testnet 1.0 launched?",
		encryptedPrice: "crqKa6gZ_hMD5ZdJYtgcBO7DJqY6nXiEZxtTESEfX7HvRy5qemXAQMUuF4buOBapQ00YcXMaFeqKSpHmAzg7CG7-PCQhhssNfXnloAm7A8mrk5LW_FZPBhq0EM4NYw1_iJtaIjD6G7Tgp3jGW5T9BC1cbIaWu9b0mzibCV7b4toOVsGGC0Zssu0GjE5PwS6n4ciUjm-eWjETE7WCntnHe3727HTBU_M=",
	},
	{
		question:       "What programming language is used by the Cosmos SDK?",
		encryptedPrice: "IE3Ny5xdnccPZ1OkgAuujs28KFHcd_tSrP_SUUDdIl8CrjhlaZJf2DzYHMlHSHzHPWeQaDnm3PNaJtI0N0Iq3slw-d_bGQquUf4HkuSLj3tB-ebqECy1_DL_wX0xOADGjBcK-BlInvNYg2XyRm85dzt1PtY8R9mlkd6LExF01G6oK1hdO7fYPjRbes2XM5AvewnlnQBzdRzE3rOOwk42H0tu",
	},
	{
		question:       "What is the name of the protocol/standard that enables communication between independent blockchains within the Cosmos network?",
		encryptedPrice: "G4N3mO8vVQz_q2YFH2zUhdcDgfYoRGF-0OLA12GE4HJ9FYWuAbhP4IiX1AM0uaMZNPSpa-xqKnP69sZFi0A-VRkrTx4BzCEX-kcR0FNs2zSa0Jk3QK78dc_dWaDqKECdyItQ6jbn0SH32AaYSpKox6hK518OC7X-gkUCOml36Os90K5RPyeblo26Qm0DuNQV0zvW6KlLttL8REf631pHGG8AQXUkukEoaIYljhWJ",
	},
	{
		question:       "What is the name of the decentralized exchange (DEX) that uses a mathematical formula to determine the price of assets based on the ratio of the available supply?",
		encryptedPrice: "PfeXaj3O8IQdQ0LxKvG0KHjnkgwQhlt-4Ce2wW45S_itrR8NIXoMc6o8nLza-GQQKNQWnkBaqKVcQHtqjyqOt16Mn86Q7Hy9IaYPbPFgqsy0f22xuLSojtgfr6z0Ur6PYlN_rST9rnaXzDvhGycxElcg5z2XIGQ-wb5DndyqU_o9aEncCTF50fXJHRtW4WwpYOc7gls1F6vOvIOuf9lqzMwa",
	},
	{
		question:       "What is the term that describes a blockchain-based exchange for trading digital assets without centralized intermediaries?",
		encryptedPrice: "ozXdDbZEk6lYR9gDBYiBX-o5Olb4ln92AYY8AtePPyuazCYnvFSbxRLdq40rNXoh5Ala1k3WWqOF4EKL7dXM_dinHUAkkBTeLghilwmPmjUYzpGC_LJrDafh6RilcPIgFMQVXSOgifSDn71pPEg258rdAKEBcgLb5Hkl7ADxt1dPmTl9-oFO6gYUC14SRZoLJSXAIn0Kehq8j9OjRjEVp5PaRt5jlBqS1dPBv5Qhe5WJpDFb",
	},
	{
		question:       "What is the name of the organization that is building and stewarding Cosmos? (abbreviation)",
		encryptedPrice: "d_eTjtu9H3a1j7d9M03ORWQN1Nq0EWvFNZulEYYK2xDNh6gYIEiun8nup84dvjSvUZUfymYS8CXKbqbJyfx6cCNiGNcjLTjrLOla-k8CCufEQV1PDiG8-Y3S6DbW6YpTt8JGbINS7oWD0_UNv8jr3ha7vRJTQ3snRURC4pJxHVC_-bwPZa0QMro3OKOhSaoj-vsUkORmO-7jRqqUKjfZKdUCEGtw8V4hYK6sJPE-lg==",
	},
	{
		question:       "What are the IRISnet community members called?",
		encryptedPrice: "xZbjWDwwkc3tWwR0X6mYuIEk9Rm8j-9oTXG5aUNmoTXLl-9CHbTEREcOtDIoZF3ADVMJGmK9z5ILvDXzajFpgPvlTq1rR_K4owTglYn8PDEeLbjWUvOwl6DtCh8bKz_tcRTLYt0nW0RAkhm1AeYDSUTs2jG0VWBgEuEddPRep5f_YgHZtKKh57YeWlkNK1l2DtmM8w7fjrWE0Lcm5uGJ3k2VEZKySfA=",
	},
	{
		question:       "What are the Cosmos community members called?",
		encryptedPrice: "_IM7T0rM4JaYjrQYuPUOsTe11ZZCetd81UgJRgDueQ2EOGN8cMwEyQfxJTTB5FjwsgSqHG4Pxdtf-FBe3SNuw2ZJFTmWkKZp06rRPiv-WGJ4x7SB8_g0ZpkoQkpO6lCQws4Na0OwMvP1eRCH3Cc7XGyfkzMwzvOdQ__uUKoYFktlJwwFjghVORyzzL9qfeSaNCnWSgLxNisQwIncrsUN-g==",
	},
	{
		question:       "What are the OmniFlix community members called?",
		encryptedPrice: "_eQ8heJ0l9-2_LQ4j-JXXCa31IIUqTyFhQI7jnc_Yxg2xkf2eYNTgzL7Epp2oPqh3UsrrvTrw7gf_y6mR7y-tynqxOKLGYJcbL5QwgtNtpy0bH1Lae1y6ybC0-28i-mojGamMGm9F12SCP9zafrAWxEvfZg3g2Z6qMYx8BEuYJuDqJdzoKNpn5nG-v_T86-9Eq5QixQD2DyNJoOkKxBcmrsuOkjUHu-8",
	},
	{
		question:       "What are the Stargaze community members called?",
		encryptedPrice: "IQBB5rQCdim283JOemmbwNy81-N6BmRaBHd53vikCfJ6pAvFqvmuUhyLGy05tb7boUMsVOLVCSGNBmbf7lXqziKaaAUmdouqo8j7Z17BwRWDg5A2xWlXBcp4ZrIIykdU3tsss0rPYFZ9yXKbLnlkzAmEMDrZRhyP7leU2l7sA_poxAhtOQ3NyIigOjexGYzyWcQRAgYS5CAnBRwqC8S7",
	},
	{
		question:       "What is the name of the newly launched state machine replication engine for the Cosmos Interchain world?",
		encryptedPrice: "P8a2P64oKCfTDW4U4aSX8OH3eyHDCbTXg3thO2nQcl0taqf0rz_zmtJPwS8k-8wDjubaSS-IOV_Hk2ZfC6nFhVFWRmFljebTt-kLiQL4I3wyoNS5NEFNLDqxlmZprFa7h41Odm1AOcwxaH_UIaHX8OXrqKgY2qEhaVHk0NhdYr16sI2-PQIDnwdRwQhr-iL7erU0jw4I6ac7UfVTpwI=",
	},
	{
		question:       "What is the delimiter for ibc class trace?",
		encryptedPrice: "vQS2tVP0k5Iut3NAMMqatUldOO5-qDeDlpgd42MxadNvzBYKaWBBtZacl6JGgvcFLKoiPwVebpGH0ui9cRywP6BoTETMJ0eJUEESGd45j0m0UOeKuOuxc8O6KIP0AsajyniQR1pEsZ__PbJenAlyRCcwb-f4C9H_Vf2vPfW_hW_QTME0S7zmndWDdY1YkGU5h_hEHnNPMFF6JEtnHfr7WmaLx5AOKwIufr-10mA2vn_VlAyXzTFAUQ==",
	},
	{
		question:       "In which year was the Cosmos Interchain Accounts feature launched?",
		encryptedPrice: "7PGJmSd9nZFUxepfGWKboMtnEZnucPEswUBpe74ciWU2JyZnKfFcaULCBBHGCnk4C6hk84WFpcihxDCTFB91i1S3GktoSYJh2nyZiwweLIiWu5YGsE_udzsZSMYje2xME3oXxFqawy9UlIalRmRO0eEKz_ic4-lQchejtaFkuC511j3_dcADja4X7he6Mv4fY-DToB1Wd6NwkjudYdvvaFi4pA3iyKzhWGTxI4C9",
	},
}

func gonQuizInteractive() {
	tryAgain := false
	var q quizQuestion
	for {
		if !tryAgain {
			q = chooseOne("Which question do you want to answer?", fullQuiz)
		}

		tryAgain = false

		a := askForString("What is the answer? (Empty answer stops the quiz)")
		if a == "" {
			break
		}

		decryptedAnswer, err := decrypt(a, q.encryptedPrice)
		if err != nil {
			println("Error decrypting answer", err.Error())
			tryAgain = true
			continue
		}

		if utf8.ValidString(decryptedAnswer) && len(strings.Split(decryptedAnswer, " ")) > 10 {
			fmt.Println("The decrypted answer:", decryptedAnswer)
			fmt.Println("If it looks like a mnemonic, the answer is correct!")
			if keepGoing := askForConfirmation("Do you want to do another one?", true); !keepGoing {
				break
			}
			continue
		}

		tryAgain = askForConfirmation("Incorrect, do you want to try again?", true)
	}
}

// Copied from https://github.com/game-of-nfts/gon-toolbox/blob/main/nft/types/aes.go
func decrypt(key string, cryptoText string) (string, error) {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
	k := generateKey(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return fmt.Sprintf("%s", cipherText), nil
}

func generateKey(key string) (genKey []byte) {
	keyBz := []byte(key)
	genKey = make([]byte, 32)
	copy(genKey, keyBz)
	for i := 32; i < len(keyBz); {
		for j := 0; j < 32 && i < len(keyBz); j, i = j+1, i+1 {
			genKey[j] ^= keyBz[i]
		}
	}
	return genKey
}
