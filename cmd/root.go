/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"io"
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"alien-invasion-cc/engine"
)

var (
	numAliens uint
	maxMoves 	uint
	mapFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alien-invasion-cc",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	 RunE: func(cmd *cobra.Command, args []string) error { 
		in, err := os.Open(mapFile)
		defer func() { _ = in.Close() }()
		if err != nil {
			return err
		}

		c := &config{
			numAliens: 	numAliens,
			maxMoves: 		maxMoves,
			in: 			in,
			out: 			cmd.OutOrStdout(),
		}
		fmt.Printf("Map File Path:%v\n", mapFile)
		fmt.Printf("Number Of Aliens:%v\n", numAliens)
		fmt.Printf("Max Moves:%v\n\n", maxMoves)



		return runEngine(cmd.Context(), c)
	 },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alien-invasion-cc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().UintVarP(&numAliens, "aliens", "n", 5, "number of aliens to be spawned")
	rootCmd.Flags().UintVarP(&maxMoves, "steps", "s", 10000, "number of maximum moves")
	rootCmd.Flags().StringVarP(&mapFile, "file", "m", "test_data/test_map", "map file path")
}

type config struct {
	numAliens, maxMoves 	uint
	in						io.ReadCloser
	out 					io.Writer
}

func runEngine(ctx context.Context, c *config) error {

	gameEngine := engine.NewEngine(
		c.numAliens,
		c.maxMoves,
		c.in,
		c.out,
	)

	return gameEngine.Run(ctx)
}
