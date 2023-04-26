#include <iostream>

template <typename T>
struct BinaryTreeNode {
    T data;
    BinaryTreeNode *leftChild;
    BinaryTreeNode *rightChild;
};

template <typename T>
class BinarySearchTree {
public:
    BinarySearchTree() {
        root = nullptr;
    }
    ~BinarySearchTree() {
        ReleaseNode(root);
    }

public:
    void inOrder() {
        inOrder(root);
    }
    void inOrder(BinaryTreeNode<T> *tNode) {
        if (tNode != nullptr) {
            inOrder(tNode->leftChild);
            std::cout << tNode->data << " ";
            inOrder(tNode->rightChild);
        }
    }

public:
    // 插入元素
    void InsertElem(const T &e) {
        InsertElem(root,e);
    }
    void InsertElem(BinaryTreeNode<T> *&tNode,const T &e) {
        // 空树 插入第一个节点
        if (tNode == nullptr) {
            tNode = new BinaryTreeNode<T>;
            tNode->data = e;
            tNode->leftChild = nullptr;
            tNode->rightChild = nullptr;
            return;
        }

        if (e > tNode->data) {
            InsertElem(tNode->rightChild,e);
        } else if (e < tNode->data) {
            InsertElem(tNode->leftChild,e);
        } else {} // no action

        return;
    }

    // 查找元素
    BinaryTreeNode<T>* SearchElem(const T &e) {
        return SearchElem(root,e);
    }
    BinaryTreeNode<T>* SearchElem(BinaryTreeNode<T> *tNode,const T &e) {
        if (tNode == nullptr) {
            return nullptr;
        }
        if (tNode->data == e) {
            return tNode;
        }

        if (e < tNode->data) {
            return SearchElem(tNode->leftChild,e);
        } else {
            return SearchElem(tNode->rightChild,e);
        }
    }

    // 删除元素
    void DeleteElem(const T& e) {
        return DeleteElem(root,e);
    }
    void DeleteElem(BinaryTreeNode<T> *&tNode,const  T &e) {
        if (tNode == nullptr) {
            return;
        }
        if (e > tNode->data) {
            DeleteElem(tNode->rightChild,e);
        } else if (e < tNode->data) {
            DeleteElem(tNode->leftChild,e);
        } else {
            // 左右子树都不为空
            if (tNode->leftChild != nullptr && tNode->rightChild != nullptr) {
                BinaryTreeNode<T> *tmpnode = tNode->leftChild;
                // 找到最右下节点
                while (tmpnode->rightChild) {
                    tmpnode = tmpnode->rightChild;
                }
                tNode->data = tmpnode->data; // 替换至待删除节点
                DeleteElem(tNode->leftChild,tmpnode->data); // 将原最右下节点删除
            } else {
                BinaryTreeNode<T> *tmpnode = tNode;
                // 更新父节点指针 将该指针指向删除节点的非空子节点
                if (tNode->leftChild == nullptr) {
                    tNode = tNode->rightChild;
                } else {
                    tNode = tNode->leftChild;
                }
                delete tmpnode; // 删除节点
            }
        }
    }

    // 查找值最大节点
    BinaryTreeNode<T>* SearchMaxValuePoint() {
        return SearchMaxValuePoint(root);
    }
    BinaryTreeNode<T>* SearchMaxValuePoint(BinaryTreeNode<T> *tNode) {
        if (tNode == nullptr) {
            return nullptr;
        }

        // 从根节点开始往右侧寻找
        BinaryTreeNode<T> *tmpnode = tNode;
        while (tmpnode->rightChild != nullptr) {
            tmpnode = tmpnode->rightChild;
        }
        return tmpnode;
    }

    // 查找值最小节点
    BinaryTreeNode<T>* SearchMinValuePoint() {
        return SearchMinValuePoint(root);
    }

    BinaryTreeNode<T>* SearchMinValuePoint(BinaryTreeNode<T> *tNode) {
        if (tNode == nullptr) {
            return nullptr;
        }

        // 从根节点开始向左侧寻找
        BinaryTreeNode<T> *tmpnode = tNode;
        while (tmpnode->leftChild != nullptr) {
            tmpnode = tmpnode->leftChild;
        }
        return tmpnode;
    }

    // 按中序遍历查找二叉树中当前节点的前趋节点
    BinaryTreeNode<T>* GetPriorPoint(BinaryTreeNode<T> *findnode) {
        if (findnode == nullptr){
            return nullptr;
        }

        BinaryTreeNode<T> *prevnode = nullptr;
        BinaryTreeNode<T> *currnode = root;
        while (currnode != nullptr) {
            if (currnode->data < findnode->data) {
                if (prevnode == nullptr) {
                    prevnode = currnode;
                } else {
                    if (prevnode->data < currnode->data) {
                        prevnode = currnode;
                    }
                }
                currnode = currnode->rightChild;
            } else if (currnode->data > findnode->data) {
                currnode = currnode->leftChild;
            } else {
                currnode = currnode->leftChild;
            }
        }

        return prevnode;
    }

    // 按中序遍历查找二叉树中当前节点的后继节点
    BinaryTreeNode<T>* GetNextPoint(BinaryTreeNode<T> *findnode) {
        if (findnode == nullptr) {
            return nullptr;
        }
        BinaryTreeNode<T> *nextnode = nullptr;
        BinaryTreeNode<T> *currnode = root;
        while (currnode != nullptr) {
            if (currnode->data > findnode->data) {
                if (nextnode == nullptr) {
                    nextnode = currnode;
                } else {
                    if (nextnode->data > currnode->data) {
                        nextnode = currnode;
                    }
                }
                currnode = currnode->leftChild;
            } else if (currnode->data < findnode->data) {
                currnode = currnode->rightChild;
            } else {
                currnode = currnode->rightChild;
            }
        }
        return nextnode;
    }


private:
    void ReleaseNode(BinaryTreeNode<T> *pnode) {
        if (pnode != nullptr) {
            ReleaseNode(pnode->leftChild);
            ReleaseNode(pnode->rightChild);
        }
        delete pnode;
    }
private:
    BinaryTreeNode<T> *root;
};

int main(void) {
    BinarySearchTree<int> tree;
    int array[] = {23,17,11,19,8,12};
    int account = sizeof(array) / sizeof(int);
    for (int i=0;i<account;i++) {
        tree.InsertElem(array[i]);
    }
    tree.inOrder();

    int val = 19;
    std::cout << std::endl;
    BinaryTreeNode<int> *tmpp;
    tmpp = tree.SearchElem(val);
    if (tmpp != nullptr) {
        std::cout << "find value: " << val << std::endl; 
    } else {
        std::cout << "miss value: " << val << std::endl;
    }

    return 0;
}