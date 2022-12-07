package imapp

type IM interface {
	Launch()
	Background()
	Foreground()
	Kill()
	Info() string
}

type App struct {
}

func (app *App) Launch() error {
	panic("please implement me")
}

func (app *App) Background() {

}

func (app *App) Foreground() {

}

func (app *App) Kill() {
	panic("please implement me")
}

func (app *App) Info() string {
	panic("please implement me")
}
