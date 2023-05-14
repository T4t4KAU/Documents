import numpy
import numpy as np

class Board(object):
    def __init__(self, **kwargs):
        self.width = int(kwargs.get('width',8))  # 棋盘宽度
        self.height = int(kwargs.get('height',8))  # 棋盘高度
        self.n_in_row = int(kwargs.get('n_in_row',5))  # 几子成线
        self.players = (1,2)  # 对弈双方编号
        self.states = {}  # 记录落子情况

    # 初始化棋盘
    def init_board(self,start_player=0):
        # 检查棋盘合法性
        if self.width < self.n_in_row or self.height < self.n_in_row:
            raise Exception(f'board width and height can not be less than {self.n_in_row}')

        # 当前player编号
        self.current_player = self.players[start_player]

        # 将可落子位初始化棋盘所有位置
        self.availables = list(range(self.width * self.height))
        self.states = {}
        self.last_move = -1

    def move_to_location(self,move):
        h = move // self.width
        w = move % self.width
        return [h, w]

    def location_to_move(self,location):
        if len(location) != 2:
            return -1

        h = location[0]
        w = location[1]
        move = h * self.width + w
        if move not in range(self.width * self.height):
            return -1

        return move

    # 走棋
    def do_move(self,move):
        self.states[move] = self.current_player
        self.availables.remove(move)
        self.current_player = (
            self.players[0] if self.current_player == self.players[1]
            else self.players[1]
        )
        self.last_move = move

    # 获取当前棋手编号
    def get_current_player(self):
        return self.current_player

    # 使用4个二值特征平面表示状态
    # 1 2平面分别表示当前player棋子位置 有子为1 无子为0
    # 3平面表示对手近一步落子 4平面表示是否为先手
    # 获取棋盘当前描述 输入给策略网络
    def current_state(self):
        square_state = np.zeros((4,self.width,self.height))
        if self.states:
            moves,players = np.array(list(zip(*self.states.items())))
            move_curr = moves[players == self.current_player]
            move_oppo = moves[players != self.current_player]
            square_state[0][move_curr // self.width, move_curr % self.height] = 1.0
            square_state[1][move_oppo // self.width, move_oppo % self.height] = 1.0
            square_state[2][self.last_move // self.width, self.last_move % self.height] = 1.0

        if len(self.states) % 2 == 0:
            square_state[3][:, :] = 1.0

        return  square_state[:, ::-1, :]

    # 判断当前棋盘有无胜方
    def has_a_winner(self):
        width = self.width
        height = self.height
        states = self.states
        n = self.n_in_row

        moved = list(set(range(width * height)) - set(self.availables))
        if len(moved) < self.n_in_row + 2:
            return False, -1

        for m in moved:
            h = m // width
            w = m % width
            player = states[m]

            if w in range(width-n+1) and len(set(states.get(i,-1) for i in range(m, m+n))) == 1:
                return True, player

            if h in range(height-n+1) and len(set(states.get(i,-1) for i in range(m, m+n*width, width))) == 1:
                return True, player

            if w in range(width-n+1) and h in range(height-n+1) and \
                    len(set(states.get(i,-1) for i in range(m, m+n*(width+1), width+1))) == 1:
                return True, player

            if w in range(n-1,width) and h in range(height-n+1) and \
                    len(set(states.get(i,-1) for i in range(m, m+n*(width-1), width-1))) == 1:
                return True, player

            return False, -1

    def game_end(self):
        win, winner = self.has_a_winner()
        if win:
            return True, winner
        elif not len(self.availables):
            return True, -1
        return False, -1


class Game(object):
    def __init__(self,board:Board):
        self.board = board

    # 自我对弈
    def start_self_play(self,player,is_shown=False,temp=1e-3):
        self.board.init_board()
        p1, p2 = self.board.players
        states, mcts_probs, current_players = [],[],[]
        while True:
            move, move_probs = player.get_action(self.board,temp=temp,return_prob=1)
            states.append(self.board.current_state())
            mcts_probs.append(move_probs)
            current_players.append(self.board.current_player)

            self.board.do_move(move)
            if is_shown:
                self.graphic(self.board,p1,p2)
            end,winner = self.board.game_end()
            if end:
                winners_z = np.zeros(len(current_players))
                if winner != -1:
                    winners_z[np.array(current_players) == winner] = 1.0
                    winners_z[np.array(current_players) != winner] = -1.0

                # 重置MCTS根节点
                player.reset_player()

                if is_shown:
                    if winner != -1:
                        print("Game over. Winner is player:",winner)
                    else:
                        print("Game end. Tie")  # 平局

                return winner, zip(states,mcts_probs,winners_z)

    # 与玩家对弈
    def start_player(self,player1,player2,start_player=0,is_shown=1):
        if start_player not in (0,1):
            raise Exception('start_player should be either 0 or 1')

        # 初始化棋盘
        self.board.init_board(start_player)
        p1, p2 = self.board.players
        player1.set_player_ind(p1)
        player2.set_player_ind(p2)
        players = {p1: player1, p2 : player2}
        if is_shown:
            self.graphic(self.board,player1.player, player2.player)
        while True:
            current_player = self.board.get_current_player()
            player_in_turn = players[current_player]
            move = player_in_turn.get_action(self.board)
            self.board.do_move(move)
            if is_shown:
                self.graphic(self.board, player1.player, player2.player)
            end, winner = self.board.game_end()
            if end:
                if is_shown:
                    if winner != -1:
                        print("Game end. Winner is",players[winner])
                    else:
                        print("Game end. Tie")

                return winner

    # 图形化展示棋盘
    def graphic(self,board:Board,player1,player2):
        width = board.width
        height = board.height

        print("Player",player1, "with X".rjust(3))
        print("Player",player2, "with O".rjust(3))
        print()

        for x in range(width):
            print("{0:8}".format(x), end='')
        print('\r\n')
        for i in range(height - 1, -1, -1):
            print("{0:4d}".format(i),end='')
            for j in range(width):
                loc = i * width + j
                p = board.states.get(loc,-1)
                if p == player1:
                    print('X'.center(8),end='')
                elif p == player2:
                    print('O'.center(8),end='')
                else:
                    print('_'.center(8),end='')
            print('\r\n\r\n')
