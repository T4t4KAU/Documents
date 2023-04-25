#include <iostream>
#include <cmath>

#define MaxSize 100

enum ECCHILDSIGN {
    E_Root,      // 树根
    E_ChildLeft, // 左孩子
    E_ChildRight // 右孩子
};

template <typename T>
struct BinaryTreeNode {
    T data;       // 数据域
    bool isValid; // 节点是否有效
};

template <typename T>
class BinaryTree {
public:
    BinaryTree() {
        for (int i=0;i<=MaxSize;i++) {
            // 初始时节点无效
            SeqBiTree[i].isValid = false;
        }
    }
    ~BinaryTree() {};

public:
    // 创建节点
    int CreateNode(int parindex,ECCHILDSIGN pointSign,const T &e);
    
    // 获取父节点下标
    int getParentIdx(int sonindex) {
        if (ifValidRangeIdx(sonindex) == false) {
            return -1;
        }
        if (SeqBiTree[sonindex].isValid == false) {
            return -1;
        }
        return int(sonindex / 2);
    }

    // 获取指定节点所在高度
    int getPointLevel(int index) {
        if (ifValidRangeIdx(index) == false) {
            return -1;
        }
        if (SeqBiTree[index].isValid == false) {
            return -1;
        }
        int level = int(log(index)/log(2)+1);
        return level;
    }

    // 获取二叉树深度
    int getLevel() {
        if (SeqBiTree[1].isValid == false) {
            return 0;
        }

        int i;
        for (i = MaxSize;i >= 1;i--) {
            // 找到最后一个有效节点
            if (SeqBiTree[i].isValid == true) {
                break;
            }
        }
        return getPointLevel(i);
    }

    // 判断是否为完全二叉树
    bool ifCompleteBT() {
        if (SeqBiTree[1].isValid == false) {
            return false;
        }
        int i;
        for (i = MaxSize;i >= 1;i--) {
            // 找到最后一个节点
            if (SeqBiTree[i].isValid == true) {
                break;
            }
        }
        for (int k = 1;k <= i;k++) {
            // 所有节点有效
            if (SeqBiTree[k].isValid == false) {
                return false;
            }
        }
        return true;
    }

    void preOrder() {
        if (SeqBiTree[1].isValid == false) {
            return;
        }
        preOrder(1);
    }

    void preOrder(int index) {
        if (ifValidRangeIdx(index) == false) {
            return;
        }
        if (SeqBiTree[index].isValid == false) {
            return;
        }
        std::cout << (char)SeqBiTree[index].data << "";
        preOrder(2 * index);
        preOrder(2 * index + 1);
    }

private:
    bool ifValidRangeIdx(int index) {
        if (index < 1 || index > MaxSize) {
            return false;
        }
        return true;
    }
private:
    BinaryTreeNode<T> SeqBiTree[MaxSize + 1];
};

template <class T>
int BinaryTree<T>::CreateNode(int parindex,ECCHILDSIGN pointSign,const T &e) {
    if (pointSign != E_Root) {
        if (ifValidRangeIdx(parindex) == false) {
            return -1;
        }
        if (SeqBiTree[parindex].isValid == false) {
            return -1;
        }
    }

    int index = -1;
    if (pointSign == E_Root) {
        index = 1; // 根节点下标为1
    } else if (pointSign == E_ChildLeft) {
        // 左孩子
        index = 2 * parindex;
        if (ifValidRangeIdx(index) == false) {
            return -1;
        }
    } else {
        // 右孩子
        index = 2 * parindex + 1;
        if (ifValidRangeIdx(index) == false) {
            return -1;
        }
    }
    SeqBiTree[index].data = e;
    // 标记该下标中有效数据
    SeqBiTree[index].isValid = true;
    return index;
}

int main(void) {
    BinaryTree<int> tree;
    int indexRoot = tree.CreateNode(-1,E_Root,'A');
    int indexNodeB = tree.CreateNode(indexRoot,E_ChildLeft,'B');
    int indexNodeC = tree.CreateNode(indexRoot,E_ChildRight,'C');

    int indexNodeD = tree.CreateNode(indexNodeB,E_ChildLeft,'D');
    int indexNodeE = tree.CreateNode(indexNodeC,E_ChildRight,'E');

    int iParentIndexE = tree.getParentIdx(indexNodeE);
    std::cout << "node E parent node index: " << iParentIndexE << std::endl;

    int iLevel = tree.getPointLevel(indexNodeD);
    std::cout << "the height of node D: " << iLevel << std::endl;

    iLevel = tree.getPointLevel(indexNodeE);
    std::cout << "the height of node E: " << iLevel << std::endl;
    std::cout << "the depth of binary tree: " << tree.getLevel() << std::endl;
    std::cout << "compelete binary tree: " << tree.ifCompleteBT() << std::endl;

    std::cout << "-----------------" << std::endl;
    std::cout << "preorder: ";
    tree.preOrder();

    return 0; 
}