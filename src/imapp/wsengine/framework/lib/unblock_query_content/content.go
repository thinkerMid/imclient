package unblockQueryContent

import (
	"ws/framework/utils"
)

var content = []string{
	"Hello please help me, I can't log in to my account",
	"Hello! My account was used by criminals, which resulted in the account being banned. Now I urge you to help unblock it, because the account contains a lot of contact information of my family and friends, which is extremely important to me. After unblocking, I will definitely take precautions here to avoid being misappropriated by criminals again. Thank you very much",
	"My account did not do anything illegal and was banned for no reason. Now my account is used to communicate with friends and family. It is very important, so please unblock my account, thank you very much",
	"Hello, my account is now unavailable. I have always cherished my account. I have not done any illegal operations. I request to unblock the account. Thank you very much",
	"Hello, I haven't been online for a long time due to work reasons. As a result, after opening it today, I found that my account has been banned. I request to unblock it. Thank you very much.",
	"There is very important data in my account, which saves my important customer information. If I can't unblock it, my work will be lost",
	"The system prompts that the account is banned, the reason for the ban is unknown, please help to unblock the account",
	"My account was banned by the system, and I didn't do anything illegal. Please help me unblock my account and tell me the reason for the ban, thank you!",
	"Hello, my number suddenly can't be used. I use it normally. There is no illegal content. Please help to investigate and unblock it as soon as possible, thank you!",
	"My account was still in normal use yesterday, and suddenly it can't be used today. Please deal with it. Now I can't contact my customers.",
	"My account has been banned. Please help to unblock it, and by the way, help to find out the reason, thank you!",
	"Suddenly found that I can't log in, but I didn't do anything illegal, why was it banned! The account has important customer information, please help to unblock it, thank you!",
	"Hello! My account system prompts that it has been banned. The reason for the ban is unknown. The account is still in use. Please help to unblock it, thank you!",
	"Hello, I don't know why my account is blocked. When I log in today, it shows that it has been banned. This account has some of my more important friends and customers. I hereby apply for unblocking, thank you",
	"Hello, my number is banned and I can't log in, please help me recover, thanks",
	"Hello, I lost my mobile phone, and my account was banned due to being used. Now I hope to unblock it and deal with the company's affairs, thank you",
	"Hello, I traveled abroad with my family. I used the VPN service for some reasons. During the process of using WhatsApp, my account was banned. There are friends I haven't seen for a long time in my account. Please unblock my account, thank you",
	"Hi brother, my account has been banned. I didn't do anything. Why do you want to ban my account? Please check the account. I haven't violated the rules? What's going on?",
	"My mother said that her whatsapp can't send information, she asked me to check it, and found that I can't log in, please help me unblock, what happened? ",
	"I need to use whatsapp to communicate with my friends, but my account is banned, please help me unblock",
	"I just sent a few more messages to the contact and was banned. Please unblock it",
	"Please unblock my account. Maybe someone is using my account, I'm not sure what happened",
	"My whatsapp account has been banned, please unblock my account, please, thank you, thank you",
	"My whatsapp account has been banned, I urgently need to contact someone, please unblock it",
	"Please unblock it. Maybe someone is using my account, I'm not sure what happened",
	"Please unblock my account, please, thank you, thank you",
	"I need to contact someone urgently now, please unblock it",
	"My whatsapp account has been banned, please unblock the account, thank you, thank you",
	"My number suddenly can't be used to unblock the account",
	"My login is abnormal, I can't log in suddenly, please check what's going on",
}

// New .
func New() string {
	size := len(content) - 1

	chooseIdx := utils.RandInt64(0, int64(size))

	return content[chooseIdx]
}
