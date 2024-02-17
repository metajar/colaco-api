package main

import (
	v1 "colaco-api/internal/api/v1"
	"colaco-api/internal/server"
	"colaco-api/internal/storage"
)

func main() {
	vendingMachine := server.NewVendingMachine(
		server.WithStorage(storage.NewMemoryStorage()),
		server.WithStartingSodas(startingSodas),
		server.WithPort("8080"),
	)
	vendingMachine.Run()
}

var startingSodas = []v1.VendingSlot{
	{
		Cost:        f322p(1.00),
		MaxQuantity: i2p(100),
		OccupiedSoda: &v1.Soda{
			Calories:    i2p(190),
			Description: s2p("An effervescent fruity experience with hints of grape and coriander."),
			Name:        s2p("Fizz"),
			OriginStory: s2p("In a quirky rooftop lab nestled in a bustling city, Dr. Effervescence, or \"Effie,\" concocted a unique beverage under the glow of a full moon. Mixing grape essence with a hint of coriander and a secret effervescent elixir, she created Fizz—an effervescent fruity experience that captured the essence of adventure in every bubble. Quickly becoming a citywide sensation for just 1 dollar US, Fizz wasn't just a drink; it was a promise of joy and a spark of excitement in every sip, born from a night of magical experimentation and destined to become legend."),
			Ounces:      f322p(16.9),
		},
		Quantity: i2p(100),
	},
	{
		Cost:        f322p(1.00),
		MaxQuantity: i2p(100),
		OccupiedSoda: &v1.Soda{
			Calories:    i2p(185),
			Description: s2p("An explosion of flavor that will knock your socks off!"),
			Name:        s2p("Pop"),
			OriginStory: s2p("In the bustling heart of a neon-lit city, amidst the clatter of creativity and the hum of innovation, Pop was born—an audacious drink that dared to challenge the mundane. Crafted by a renegade chef known only as \"The Flavor Maverick\" in a clandestine urban kitchen, Pop emerged as a defiant explosion of flavors, destined to jolt the taste buds and electrify the senses. With a secret blend that promised an adventure in every gulp, Pop became the talk of the town, a beacon for thrill-seekers and flavor chasers alike. Priced at just 1 dollar US and with only 100 bottles available to vend, it wasn't just a beverage—it was a treasure hunt for the palate, a limited-edition experience that promised to knock your socks off with every effervescent sip. Pop wasn't merely a drink; it was a revolution in a bottle, waiting to unleash an explosion of flavor with the power to transform the ordinary into the extraordinary."),
			Ounces:      f322p(16.9),
		},
		Quantity: i2p(100),
	},
	{
		Cost:        f322p(1.00),
		MaxQuantity: i2p(200),
		OccupiedSoda: &v1.Soda{
			Calories:    i2p(225),
			Description: s2p("A basic no nonsense cola that is the perfect pick me up for any occasion."),
			Name:        s2p("Cola"),
			OriginStory: s2p("In the heart of a small, bustling town where traditions meld seamlessly with the pulse of modern life, Cola was born—a straightforward, no-nonsense drink crafted for the soul of simplicity. In a world brimming with complex flavors and endless choices, a local beverage artisan, affectionately known as \"Old Joe,\" decided it was time to return to the basics. With a timeless recipe, he created Cola, a classic cola that became an instant favorite. Its familiar taste was like a comforting embrace, the perfect pick-me-up for any occasion. Priced at just 1 dollar US and with a generous 200 bottles available to vend, Cola captured the essence of what it means to enjoy the simpler things in life. It wasn't trying to be a fleeting trend or a collector's craze; it was the essence of reliability and refreshment, a testament to the power of keeping things simple and sweet. Cola became more than just a drink; it was a staple, a reminder that sometimes, the most basic pleasures are the ones that truly satisfy."),
			Ounces:      f322p(16.9),
		},
		Quantity: i2p(200),
	},
	{
		Cost:        f322p(1.00),
		MaxQuantity: i2p(150),
		OccupiedSoda: &v1.Soda{
			Calories:    i2p(356),
			Description: s2p("Not for the faint of heart.  So flavorful and so invigorating, it should probably be illegal."),
			Name:        s2p("Mega Pop"),
			OriginStory: s2p("In the shadowy corners of a city that never sleeps, where the thrill of the forbidden dances on the tongues of the daring, Mega Pop was concocted. This elixir, born from the genius of an underground flavor wizard known only as \"The Alchemist,\" was a defiant act against the mundane. Mega Pop was not just a drink; it was a rebellion in a bottle, bursting with flavors so bold and invigorating they bordered on the edge of legality. Crafted for those who seek the extreme, its recipe was whispered to be a fusion of exotic ingredients from hidden corners of the world, each sip a testament to the audacity of its creation. Priced at a mere 1 dollar US and limited to only 50 bottles in circulation, Mega Pop became the urban legend everyone had to taste to believe. It wasn't for the faint of heart—it was a dare, a challenge, a thrilling ride for the palate that promised an experience as unrivaled as it was unforgettable. Mega Pop: a beverage so intense, it flirted with the limits of the law, offering a taste of the wild side to those brave enough to take the plunge."),
			Ounces:      f322p(16.9),
		},
		Quantity: i2p(50),
	},
}

func s2p(s string) *string {
	return &s
}

func i2p(i int) *int {
	return &i
}

func f322p(f float32) *float32 {
	return &f
}
