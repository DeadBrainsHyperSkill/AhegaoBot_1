package discordmarket

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Provides common functionality to all reaction subfunctions.
func (b *Bot) reactionBase(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var userId = i.Member.User.ID

	var dial = i.ApplicationCommandData().Options[0].IntValue()

	var userMention = i.Member.User.Mention()

	var title string

	var color int

	var randImage string

	var randDescription string

	switch dial {
	case 25:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/938421483057860658/bleh_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938421481224962058/bleh_double1.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938421481828925480/bleh_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938421482596495380/bleh_double3.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v задирает %v, кому-то не поздоровиться.**", userMention, subjectMention),
				fmt.Sprintf("**%v понимает, что %v решил поиздеваться над ним.**", subjectMention, userMention),
				fmt.Sprintf("**%v искривляет лицо перед %v, опять он за своё!**", userMention, subjectMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/938421452506537984/bleh.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938421451009187940/bleh1.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938421451420237844/bleh2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938421452020011038/bleh3.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**%v тот ещё забияка, он дразнит всех.**", userMention),
				fmt.Sprintf("**%v корчит рожу, прямо перед всеми!**", userMention),
				fmt.Sprintf("**Кто задирает и всех обижает? Кроме %v, никто.**", userMention),
			})
		}

		title = "Реакция: Дразнить"
		color = 16737894

	case 26:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/909119404791451658/bite_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909119383933161472/bite_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909119386672066620/bite_double3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909119390753099817/bite_double4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909119398059601950/bite_double5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909119398143463464/bite_double6.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938417205383667732/bite_double7.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938417205656309831/bite_double8.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938417205085888592/bite_double9.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v подвергся укусу %v, было больно.**", subjectMention, userMention),
				fmt.Sprintf("**%v хочет урвать кусочек от %v**", userMention, subjectMention),
				fmt.Sprintf("**%v голоден, и решил покусать %v**", userMention, subjectMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/909119379302645810/bite.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909120671748075580/bite2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909119409627467776/bite3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909119412710281226/bite4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909119370192650271/bite5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/909119376148541440/bite6.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940943694750707742/bite7.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**%v делает очередной кусь.**", userMention),
				fmt.Sprintf("**%v собрался оторвать лакомый кусок.**", userMention),
				fmt.Sprintf("**Хруст-хруст, кажется %v кусает кого-то или что-то?**", userMention),
			})
		}

		title = "Реакция: Кусь"
		color = 16711680

	case 27:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/941028321058250752/kneel.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/941028321699954708/kneel2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/941028322605936740/kneel3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/941028325164453958/kneel4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/941028319925792858/kneel5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/941028320294866994/kneel6.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v беспрекословно встаёт на колени перед %v**", userMention, subjectMention),
				fmt.Sprintf("**Исключительная преданность %v перед %v заметна сразу.**", userMention, subjectMention),
				fmt.Sprintf("**Под натиском %v, даже %v вынужден склонить колени.**", subjectMention, userMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/941028765167935528/kneel_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/941028762470989874/kneel_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/941028762869440542/kneel_double3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/941028763506995220/kneel_double4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/941028764907880538/kneel_double5.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**%v встаёт на колени перед всеми.**", userMention),
				fmt.Sprintf("**А кто это у нас на коленках? Оказывается %v отдаёт честь каждому здесь.**", userMention),
				fmt.Sprintf("**Эмм... %v, ты чего колени пачкаешь?**", userMention),
			})
		}

		title = "Реакция: На колени"
		color = 8357810

	case 28:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940951801258201118/retribution_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940951797416214658/retribution_double1.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940951798976499732/retribution_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940951799853117510/retribution_double3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940951800721334322/retribution_double4.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v решил применить тяжёлую артиллерию. %v сейчас будет больно.**", userMention, subjectMention),
				fmt.Sprintf("**\"%v, если ты веришь в богов, то начинай им молиться\" - %v**", subjectMention, userMention),
				fmt.Sprintf("**Цель захвачена — %v. Ожидаю подтверждения от %v... ААААГООНЬ!!!**", subjectMention, userMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940950478492143697/retribution.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940950478999650355/retribution2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940950479721078834/retribution3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940950480555741194/retribution4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940950481214263306/retribution5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940950476650868767/retribution7.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940950477200314388/retribution8.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940950477900755036/retribution9.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**Случайный выстрел от %v - случайная пуля для...**", userMention),
				fmt.Sprintf("**Советуем всем пригнуться - %v достал свой ствол!**", userMention),
				fmt.Sprintf("**Кто же сегодня будет мишенью для %v?**", userMention),
			})
		}

		title = "Реакция: Оружие Возмездия"
		color = 16187136

	case 29:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940635418943168592/slag_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940635415629676544/slag_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940635416464334848/slag_double3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940635417353543751/slag_double4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940635418276266015/slag_double5.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v смотрит на %v, как на мусор.**", userMention, subjectMention),
				fmt.Sprintf("**\"Какой же ты жалкий\" - говорит %v про %v**", userMention, subjectMention),
				fmt.Sprintf("**%v не больше чем тля, по мнению %v**", subjectMention, userMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940635360122245180/slag.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940635361061777408/slag2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940635361904840704/slag3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940635362609487943/slag4.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**%v считает каждого ничтожеством.**", userMention),
				fmt.Sprintf("**%v смотрит на всех, как на мусор.**", userMention),
				fmt.Sprintf("**Каждый сидящий здесь - грязь под ногтями, считает %v**", userMention),
			})
		}

		title = "Реакция: Отброс"
		color = 10253096

	case 30:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940189156431257640/lick_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940189154841600010/lick_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940189155147776040/lick_double3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940189155701436416/lick_double4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940189155953111070/lick_double5.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v всё тщательно вылизал у %v. Больше не осталось...**", userMention, subjectMention),
				fmt.Sprintf("**%v старается угодить %v, выжимая все соки изнутри.**", userMention, subjectMention),
				fmt.Sprintf("**%v не ожидал такого смачного подарка от %v**", subjectMention, userMention),
				fmt.Sprintf("**%v уже на грани... А %v уже просто не лезет.**", subjectMention, userMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940189107391438888/lick.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940189107760562196/lick2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940189109169815622/lick4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940189109467619409/lick5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940189109840916551/lick6.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**%v удалил себе ребра и отсосал себе сам.**", userMention),
				fmt.Sprintf("**Ммм... %v знает толк в нежно-сладком отлизе.**", userMention),
				fmt.Sprintf("**Отлизать так, что не окажется сухого места... Так может только %v**", userMention),
			})
		}

		title = "Реакция: Отлизать"
		color = 16711880

	case 31:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/938432107519303790/hi_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938432106277789696/hi_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938432106542018630/hi_double3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938432106844004362/hi_double4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938432107125014578/hi_double5.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v тебе привет от %v**", subjectMention, userMention),
				fmt.Sprintf("**%v передает приветик %v**", userMention, subjectMention),
				fmt.Sprintf("**%v горит желанием поздороваться с %v**", userMention, subjectMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/938430061877534750/hi.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938430062393446400/hi2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938430062867410974/hi3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938430063223914586/hi4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938430063462977636/hi5.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**%v приветствует всех.**", userMention),
				fmt.Sprintf("**%v передаёт каждому привет!**", userMention),
				fmt.Sprintf("**А кто это машет всем ручкой? %v, не ты ли это?**", userMention),
			})
		}

		title = "Реакция: Приветик"
		color = 9820854

	case 32:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940652593535942666/sleep_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940652590633484308/sleep_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940652591120015390/sleep_double3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940652591799500860/sleep_double4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940652593250725918/sleep_double6.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v тихонько укладывает %v в кроватку.**", userMention, subjectMention),
				fmt.Sprintf("**%v желает сладких снов %v**", userMention, subjectMention),
				fmt.Sprintf("**%v уже в постели, и %v вслед за ним.**", subjectMention, userMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940652152651669514/sleep.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940652151217225818/sleep2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940652151586316318/sleep3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940652151854743552/sleep5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940652152269975552/sleep6.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**%v спит без задних ног.**", userMention),
				fmt.Sprintf("**Кто-то отключился, %v не беспокоить!**", userMention),
				fmt.Sprintf("**%v отдыхает после тяжелого дня.**", userMention),
			})
		}

		title = "Реакция: Спать"
		color = 9095131

	case 33:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940667076199194634/suplex_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940667072856354906/suplex_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940667073854578758/suplex_double4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940667074617958480/suplex_double5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940667075494555678/suplex_double6.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v кидает на прогиб %v**", userMention, subjectMention),
				fmt.Sprintf("**Через прогиб - последнее, что услышал %v от %v**", subjectMention, userMention),
				fmt.Sprintf("**%v сейчас почувствует на себе бросок от %v**", subjectMention, userMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940666153846575154/suplex.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/847896613417779230/suplex1.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940666154467352596/suplex2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940666155130048522/suplex3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940666155985694750/suplex4.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**%v кидает всех и вся на прогиб.**", userMention),
				fmt.Sprintf("**%v хочет всем показать свой суплекс.**", userMention),
				fmt.Sprintf("**Сейчас будет техника коронного броска от %v**", userMention),
			})
		}

		title = "Реакция: Суплекс"
		color = 2080588

	case 34:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/938436075389599754/pants_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938436076043927592/pants_double3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938436076694032434/pants_double4.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**\"Хочешь увидеть мои трусики?\" - спрашивает %v у %v**", userMention, subjectMention),
				fmt.Sprintf("**%v сейчас узрит всю красоту нижнего белья %v**", subjectMention, userMention),
				fmt.Sprintf("**%v снимает свои трусики прямо перед %v**", userMention, subjectMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/938435633217679360/pants.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938435634035572746/pants2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938435634605985832/pants3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/938435635306438657/pants4.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**%v снимает свои трусики у всех на виду!**", userMention),
				fmt.Sprintf("**Кажется сейчас мы увидем пантсу %v**", userMention),
				fmt.Sprintf("**%v хочет всем показать свои трусики.**", userMention),
			})
		}

		title = "Реакция: Трусики"
		color = 14452096

	case 35:

		if len(i.ApplicationCommandData().Options) != 1 {

			var subjectMention = i.ApplicationCommandData().Options[1].UserValue(s).Mention()

			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940619459972984842/shrug_double.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940619458299457626/shrug_double2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940619458819543150/shrug_double3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940619459335430264/shrug_double4.gif",
			})

			randDescription = randSlice([]string{
				fmt.Sprintf("**%v не понимает, о чём говорит %v**", userMention, subjectMention),
				fmt.Sprintf("**%v рассказывает, но %v может только пожать плечами.**", subjectMention, userMention),
				fmt.Sprintf("**%v слушает %v и не понимает.**", userMention, subjectMention),
			})
		} else {
			randImage = randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/940619442486923355/shrug.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940619442822479872/shrug2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940619441073434634/shrug3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/940619441895514123/shrug4.gif",
			})
			randDescription = randSlice([]string{
				fmt.Sprintf("**Откуда %v знает это?**", userMention),
				fmt.Sprintf("**%v не в курсе дела.**", userMention),
				fmt.Sprintf("**%v не знает и только пожимает плечами.**", userMention),
			})
		}

		title = "Реакция: Хз"
		color = 14452096

	}
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Color:       color,
		Description: randDescription,
		Image:       &discordgo.MessageEmbedImage{URL: randImage},
	}

	hasProduct, err := b.MarketData.HasProduct(userId, int(dial))
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to determine whether user has a product")
		return
	}

	if hasProduct {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed}}})
		return
	}

	isEnoughCurrency, balance, err := b.isEnoughCurrency(userId, b.Cfg.ReactionPrice)
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user balance")
		return
	}

	if !isEnoughCurrency {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Content: fmt.Sprintf("Одноразовое применение реакции стоит `%v` %v\n\nОднако, на твоем балансе `%v` %v\n",
					balance, b.Cfg.Currency, b.Cfg.ReactionPrice, b.Cfg.Currency),
			},
		})
		return
	}

	err = b.MarketData.BalanceUpdate(userId, -b.Cfg.ReactionPrice)
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update user's balance")
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed}}})
}
