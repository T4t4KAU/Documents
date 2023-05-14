from __future__ import print_function
import pickle
from game import Board, Game
from mcts import MCTSPlayer
from network import PolicyValueNet


class Human:
    def __init__(self):
        self.player = None

    def set_player_ind(self, p):
        self.player = p

    def get_action(self, board):
        try:
            location = input("> ")
            if isinstance(location,str):
                location = [int(n,10) for n in location(",")]
        except:
            move = -1

        if move == -1 or move not in board.availables:
            print("invalid move")

        return move

    def __str__(self):
        return "Human {}".format(self.player)


def run():
    n = 5
    width, height = 8,8
    model_file = ""
    try:
        board = Board(width=width,height=height,n_in_row=n)
        game = Game(board)

        try:
            policy_param = pickle.load(open(model_file),'rb')
        except:
            policy_param = pickle.load(open(model_file),'rb')

        best_policy = PolicyValueNet(width,height,policy_param)
        mcts_player = MCTSPlayer(best_policy.policy_value_fn, c_puct=5, n_playout=400)
        human = Human()
        game.start_player(human,mcts_player,start_player=1,is_shown=True)

    except KeyboardInterrupt:
        print('\n\rquit')
